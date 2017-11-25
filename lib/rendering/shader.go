package rendering

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/stojg/graphics/lib/components"
	"github.com/stojg/graphics/lib/physics"
)

var loadedShaders = make(map[string]*ShaderResource)
var shaderInUse uint32 = 2 ^ 32 - 1

func NewShader(fileName string) *Shader {

	s := &Shader{
		filename: fileName,
	}

	if oldResource, ok := loadedShaders[fileName]; ok {
		s.resource = oldResource
		s.resource.AddReference()
	} else {
		s.resource = NewShaderResource()
	}

	vertexShaderText := s.loadShader(fileName + ".vs")
	fragmentShaderText := s.loadShader(fileName + ".fs")

	vShader := s.addVertexShader(vertexShaderText)
	defer s.cleanUp(vShader)

	fShader := s.addFragmentShader(fragmentShaderText)
	defer s.cleanUp(fShader)

	s.CompileShader()

	s.AddAllUniforms(vertexShaderText)
	s.AddAllUniforms(fragmentShaderText)

	loadedShaders[fileName] = s.resource

	return s
}

type Shader struct {
	filename string
	resource *ShaderResource
}

func (s *Shader) UpdateUniforms(transform *physics.Transform, mat components.Material, engine components.RenderingEngine) {

	for i, name := range s.resource.uniformNames {
		uniformType := s.resource.uniformTypes[i]

		// x_ signifies that this uniform is not set by the shader directly, ie a hack
		if strings.Index(name, "x_") == 0 {
			if uniformType == "sampler2D" {
				samplerSlot := engine.GetSamplerSlot(name)
				engine.GetTexture(name).Bind(samplerSlot)
				gl.Uniform1i(s.resource.uniforms[name], int32(samplerSlot))
			} else {
				s.SetUniform(name, engine.GetVector3f(name))
			}
			continue
		}

		if uniformType == "sampler2D" {
			samplerSlot := engine.GetSamplerSlot(name)
			mat.Texture(name).Bind(samplerSlot)
			s.SetUniform(name, int32(samplerSlot))
			continue
		}

		switch name {
		case "lightMVP":
			s.SetUniform(name, engine.GetActiveLight().ViewProjection().Mul4(transform.Transformation()))
		case "projection":
			s.SetUniform(name, engine.GetMainCamera().GetProjection())
		case "model":
			s.SetUniform(name, transform.Transformation())
		case "view":
			s.SetUniform(name, engine.GetMainCamera().GetView())
		case "viewPos":
			s.SetUniform(name, engine.GetMainCamera().Transform().Pos())
		case "directionalLight":
			s.setUniformDirectionalLight(name, engine.GetActiveLight().(components.DirectionalLight))
		case "pointLight":
			s.setUniformPointLight(name, engine.GetActiveLight().(components.PointLight))
		case "spotLight":
			s.setUniformSpotLight(name, engine.GetActiveLight().(components.Spotlight))

		default:
			fmt.Printf("Shader.UpdateUniforms: unknow uniform %s\n", name)
		}

	}
}

func (s *Shader) SetUniform(uniformName string, value interface{}) {
	loc, ok := s.resource.uniforms[uniformName]
	if !ok {
		panic(fmt.Sprintf("no shader location found for uniform: '%s' in shader '%s'", uniformName, s.filename))
	}
	switch v := value.(type) {
	case float32:
		gl.Uniform1f(loc, v)
	case int32:
		gl.Uniform1i(loc, v)
	case mgl32.Mat4:
		gl.UniformMatrix4fv(loc, 1, false, &v[0])
	case mgl32.Vec3:
		gl.Uniform3fv(loc, 1, &v[0])
	default:
		panic(fmt.Sprintf("unknown uniform type for '%s'", uniformName))
	}
}

func (s *Shader) setUniformBaseLight(uniformName string, baseLight components.Light) {
	s.SetUniform(uniformName+".color", baseLight.Color())
}

func (s *Shader) setUniformDirectionalLight(uniformName string, directional components.DirectionalLight) {
	s.setUniformBaseLight(uniformName+".base", directional)
	s.SetUniform(uniformName+".direction", directional.Direction())
}

func (s *Shader) setUniformPointLight(uniformName string, pointLight components.PointLight) {
	s.setUniformBaseLight(uniformName+".base", pointLight)
	s.SetUniform(uniformName+".position", pointLight.Position())
	s.SetUniform(uniformName+".atten.constant", pointLight.Constant())
	s.SetUniform(uniformName+".atten.linear", pointLight.Linear())
	s.SetUniform(uniformName+".atten.exponent", pointLight.Exponent())
}

func (s *Shader) setUniformSpotLight(uniformName string, spotLight components.Spotlight) {
	s.setUniformPointLight(uniformName+".pointLight", spotLight)
	s.SetUniform(uniformName+".direction", spotLight.Direction())
	s.SetUniform(uniformName+".cutoff", spotLight.Cutoff())
}

func (s *Shader) Bind() {
	if shaderInUse != s.resource.Program {
		shaderInUse = s.resource.Program
		gl.UseProgram(s.resource.Program)
	}
}

func (s *Shader) loadShader(filepath string) string {
	b, err := ioutil.ReadFile("./res/shaders/" + filepath)
	if err != nil {
		panic(err)
	}
	return string(b) + "\x00"
}

func (s *Shader) addVertexShader(shader string) uint32 {
	return s.addProgram(shader, gl.VERTEX_SHADER)
}

func (s *Shader) addGeometryShader(shader string) uint32 {
	return s.addProgram(shader, gl.GEOMETRY_SHADER)
}

func (s *Shader) addFragmentShader(shader string) uint32 {
	return s.addProgram(shader, gl.FRAGMENT_SHADER)
}

func (s *Shader) AddAllUniforms(shaderText string) {

	structs := s.FindUniformStructs(shaderText)

	r := regexp.MustCompile(`uniform\s*(\S*)\s(\S*);`)

	for _, line := range strings.Split(shaderText, "\n") {
		t := r.FindAllStringSubmatch(line, -1)
		for _, i := range t {
			if len(i) == 3 {
				s.resource.AddUniformName(i[2])
				s.resource.AdduniformType(i[1])
				s.AddUniform(i[1], i[2], structs)
			}
		}
	}
}

func (s *Shader) AddUniform(glType, name string, structs map[string][]GLSLStruct) {
	structComponents, ok := structs[glType]
	if ok {
		for _, v := range structComponents {
			s.AddUniform(v.stype, name+"."+v.name, structs)
		}
		return
	}
	t := gl.GetUniformLocation(s.resource.Program, gl.Str(name+"\x00"))
	if t < 0 {
		fmt.Printf("uniform '%s' seems to not be used for shader '%s'\n", name, s.filename)
	}
	s.resource.uniforms[name] = t
}

type GLSLStruct struct {
	name  string
	stype string
}

func (s Shader) FindUniformStructs(shaderText string) map[string][]GLSLStruct {
	result := make(map[string][]GLSLStruct)
	var re = regexp.MustCompile(`(?s)struct\s*(\w*)\s+{\s([^}]*)};`)
	for _, match := range re.FindAllStringSubmatch(shaderText, -1) {
		structName := match[1]
		content := match[2]
		var inner = regexp.MustCompile(`(?s)\s*(\w*)\s(\w*);`)
		var properties []GLSLStruct
		for _, innerMatch := range inner.FindAllStringSubmatch(content, -1) {
			properties = append(properties, GLSLStruct{
				name:  innerMatch[2],
				stype: innerMatch[1],
			})
		}
		result[structName] = properties
	}
	return result
}

func (s *Shader) CompileShader() {

	gl.LinkProgram(s.resource.Program)

	var status int32
	gl.GetProgramiv(s.resource.Program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(s.resource.Program, gl.INFO_LOG_LENGTH, &logLength)

		l := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(s.resource.Program, logLength, nil, gl.Str(l))

		panic(fmt.Errorf("failed to link Program[%d]: %v", s.resource.Program, l))
	}
}

func (s *Shader) cleanUp(shader uint32) {
	gl.DetachShader(s.resource.Program, shader)
	gl.DeleteShader(shader)
}

func (s *Shader) addProgram(text string, shaderType uint32) uint32 {

	shader := gl.CreateShader(shaderType)

	if shader == 0 {
		fmt.Println("Shader creation failed: Could not find valid memory location when adding shader")
		os.Exit(1)
	}

	shaderSource, free := gl.Strs(text)
	gl.ShaderSource(shader, 1, shaderSource, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		printInfoLog(shader, text)
		os.Exit(1)
	}
	gl.AttachShader(s.resource.Program, shader)
	return shader
}

func printInfoLog(shader uint32, shaderText string) {
	var logLength int32
	gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)
	infoLog := strings.Repeat("\x00", int(logLength+1))
	gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(infoLog))
	fmt.Printf("Shader compilation failed:\n%s------------\n", infoLog)
	for i, line := range strings.Split(shaderText, "\n") {
		fmt.Printf("%d: %s\n", i, line)
	}
}
