package rendering

import (
	"bufio"
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

	vertexShaderText, err := s.loadShader(fileName + ".vert")
	if err != nil {
		fmt.Printf("error loading shader '%s': %v\n", fileName+".vert", err)
		os.Exit(1)
	}
	fragmentShaderText, err := s.loadShader(fileName + ".frag")
	if err != nil {
		fmt.Printf("error loading shader '%s': %v\n", fileName+".frag", err)
		os.Exit(1)
	}

	vShader := s.addVertexShader(vertexShaderText)
	defer s.cleanUp(vShader)

	fShader := s.addFragmentShader(fragmentShaderText)
	defer s.cleanUp(fShader)

	s.compile()

	s.addAllUniforms(vertexShaderText)
	s.addAllUniforms(fragmentShaderText)

	loadedShaders[fileName] = s.resource

	return s
}

type Shader struct {
	filename string
	resource *ShaderResource
}

func (s *Shader) Bind() {
	if shaderInUse != s.resource.Program {
		shaderInUse = s.resource.Program
		gl.UseProgram(s.resource.Program)
	}
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
				switch uniformType {
				case "bool":
					s.setUniform(name, engine.GetInteger(name))
				case "vec3":
					s.setUniform(name, engine.GetVector3f(name))
				case "float":
					s.setUniform(name, engine.GetFloat(name))
				default:
					panic(uniformType)
				}
			}
			continue
		}

		if uniformType == "sampler2D" {
			samplerSlot := engine.GetSamplerSlot(name)
			mat.Texture(name).Bind(samplerSlot)
			s.setUniform(name, int32(samplerSlot))
			continue
		}

		switch name {
		case "lightMVP":
			s.setUniform(name, engine.GetActiveLight().ViewProjection().Mul4(transform.Transformation()))
		case "projection":
			s.setUniform(name, engine.GetMainCamera().GetProjection())
		case "model":
			s.setUniform(name, transform.Transformation())
		case "view":
			s.setUniform(name, engine.GetMainCamera().GetView())
		case "viewPos":
			s.setUniform(name, engine.GetMainCamera().Pos())
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

func (s *Shader) loadShader(filepath string) (string, error) {
	b, err := ioutil.ReadFile("./res/shaders/" + filepath)
	if err != nil {
		return "", err
	}

	return s.addIncludes(string(b))
}

func (s *Shader) addIncludes(shaderText string) (string, error) {

	var re = regexp.MustCompile(`#include\s"([^"]*)"`)

	var result string

	scanner := bufio.NewScanner(strings.NewReader(shaderText))
	for scanner.Scan() {
		line := scanner.Text()

		includes := re.FindAllStringSubmatch(line, -1)
		if len(includes) > 0 {
			for _, match := range includes {
				text, err := s.loadShader(match[1])
				if err != nil {
					return "", err
				}
				result += text
			}
		} else {
			result += fmt.Sprintf("%s\n", line)
		}

	}

	return result, nil
}

func (s *Shader) setUniform(uniformName string, value interface{}) {
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
	s.setUniform(uniformName+".color", baseLight.Color())
}

func (s *Shader) setUniformDirectionalLight(uniformName string, directional components.DirectionalLight) {
	s.setUniformBaseLight(uniformName+".base", directional)
	s.setUniform(uniformName+".direction", directional.Direction())
}

func (s *Shader) setUniformPointLight(uniformName string, pointLight components.PointLight) {
	s.setUniformBaseLight(uniformName+".base", pointLight)
	s.setUniform(uniformName+".position", pointLight.Position())
	s.setUniform(uniformName+".atten.constant", pointLight.Constant())
	s.setUniform(uniformName+".atten.linear", pointLight.Linear())
	s.setUniform(uniformName+".atten.exponent", pointLight.Exponent())
}

func (s *Shader) setUniformSpotLight(uniformName string, spotLight components.Spotlight) {
	s.setUniformPointLight(uniformName+".pointLight", spotLight)
	s.setUniform(uniformName+".direction", spotLight.Direction())
	s.setUniform(uniformName+".cutoff", spotLight.Cutoff())
}

func (s *Shader) addVertexShader(shader string) uint32 {
	return s.createProgram(shader, gl.VERTEX_SHADER)
}

func (s *Shader) addGeometryShader(shader string) uint32 {
	return s.createProgram(shader, gl.GEOMETRY_SHADER)
}

func (s *Shader) addFragmentShader(shader string) uint32 {
	return s.createProgram(shader, gl.FRAGMENT_SHADER)
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

func (s *Shader) addAllUniforms(shaderText string) {

	uniformStructs := s.findUniformStructs(shaderText)

	r := regexp.MustCompile(`uniform\s*(\S*)\s(\S*);`)

	for _, line := range strings.Split(shaderText, "\n") {
		t := r.FindAllStringSubmatch(line, -1)
		for _, i := range t {
			if len(i) == 3 {
				s.resource.AddUniformName(i[2])
				s.resource.AdduniformType(i[1])
				s.AddUniform(i[1], i[2], uniformStructs)
			}
		}
	}
}

type GLSLStruct struct {
	name  string
	stype string
}

func (s Shader) findUniformStructs(shaderText string) map[string][]GLSLStruct {
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

func (s *Shader) compile() {

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

func (s *Shader) createProgram(text string, shaderType uint32) uint32 {

	shader := gl.CreateShader(shaderType)

	if shader == 0 {
		fmt.Println("Shader creation failed: Could not find valid memory location when adding shader")
		os.Exit(1)
	}

	shaderSource, free := gl.Strs(text + "\x00")
	gl.ShaderSource(shader, 1, shaderSource, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		printInfoLog(shader, s.filename, text)
		os.Exit(1)
	}
	gl.AttachShader(s.resource.Program, shader)
	return shader
}

func printInfoLog(shader uint32, filename, shaderText string) {
	var logLength int32
	gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)
	infoLog := strings.Repeat("\x00", int(logLength+1))
	gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(infoLog))
	fmt.Printf("Shader compilation failed (%s):\n%s------------\n", filename, infoLog)
	for i, line := range strings.Split(shaderText, "\n") {
		fmt.Printf("%d: %s\n", i, line)
	}
}
