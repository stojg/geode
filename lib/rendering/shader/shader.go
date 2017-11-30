package shader

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/stojg/graphics/lib/components"
	"github.com/stojg/graphics/lib/physics"
)

var loadedShaders = make(map[string]*ShaderResource)
var shaderInUse uint32 = 2 ^ 32 - 1

var shaderFolder = "./res/shaders/"

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

	s.findAllUniforms(vertexShaderText)
	s.findAllUniforms(fragmentShaderText)

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

		// x_ means that the uniform is a engine uniform and should be fetch from the render engines storage
		if strings.Index(name, "x_") == 0 {
			if uniformType == "sampler2D" {
				samplerSlot := engine.SamplerSlot(name)
				engine.Texture(name).Bind(samplerSlot)
				gl.Uniform1i(s.resource.uniforms[name], int32(samplerSlot))
			} else {
				switch uniformType {
				case "bool":
					s.updateUniform(name, engine.Integer(name))
				case "vec3":
					s.updateUniform(name, engine.Vector3f(name))
				case "float":
					s.updateUniform(name, engine.Float(name))
				default:
					panic(uniformType)
				}
			}
			continue
		}

		if uniformType == "sampler2D" {
			samplerSlot := engine.SamplerSlot(name)
			mat.Texture(name).Bind(samplerSlot)
			s.updateUniform(name, int32(samplerSlot))
			continue
		}

		switch name {
		case "projection":
			s.updateUniform(name, engine.MainCamera().Projection())
		case "model":
			s.updateUniform(name, transform.Transformation())
		case "view":
			s.updateUniform(name, engine.MainCamera().View())
		case "lightMVP":
			s.updateUniform(name, engine.ActiveLight().ViewProjection().Mul4(transform.Transformation()))
		case "directionalLight":
			s.updateUniformDirectionalLight(name, engine.ActiveLight().(components.DirectionalLight))
		case "pointLight":
			s.updateUniformPointLight(name, engine.ActiveLight().(components.PointLight))
		case "spotLight":
			s.updateUniformSpotLight(name, engine.ActiveLight().(components.Spotlight))
		default:
			fmt.Printf("Shader.UpdateUniforms: unknow uniform %s\n", name)
		}

	}
}

func (s *Shader) loadShader(filepath string) (string, error) {
	b, err := ioutil.ReadFile(shaderFolder + filepath)
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

		if len(includes) == 0 {
			result += fmt.Sprintf("%s\n", line)
			continue
		}

		for _, include := range includes {
			text, err := s.loadShader(include[1])
			if err != nil {
				return "", err
			}
			result += text
			result += fmt.Sprintf("// ^ %s\n", include[1])
		}
	}

	return result, nil
}

func (s *Shader) updateUniform(uniformName string, value interface{}) {
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
	s.updateUniform(uniformName+".color", baseLight.Color())
	s.updateUniform(uniformName+".maxDistance", baseLight.MaxDistance())
}

func (s *Shader) updateUniformDirectionalLight(uniformName string, directional components.DirectionalLight) {
	s.setUniformBaseLight(uniformName+".base", directional)
	s.updateUniform(uniformName+".direction", directional.Direction())
}

func (s *Shader) updateUniformPointLight(uniformName string, pointLight components.PointLight) {
	s.setUniformBaseLight(uniformName+".base", pointLight)
	s.updateUniform(uniformName+".position", pointLight.Position())
	s.updateUniform(uniformName+".atten.constant", pointLight.Constant())
	s.updateUniform(uniformName+".atten.linear", pointLight.Linear())
	s.updateUniform(uniformName+".atten.exponent", pointLight.Exponent())
}

func (s *Shader) updateUniformSpotLight(uniformName string, spotLight components.Spotlight) {
	s.updateUniformPointLight(uniformName+".pointLight", spotLight)
	s.updateUniform(uniformName+".direction", spotLight.Direction())
	s.updateUniform(uniformName+".cutoff", spotLight.Cutoff())
}

func (s *Shader) addVertexShader(shader string) uint32 {
	return s.createProgram(shader, gl.VERTEX_SHADER)
}

func (s *Shader) addFragmentShader(shader string) uint32 {
	return s.createProgram(shader, gl.FRAGMENT_SHADER)
}

func (s *Shader) findAllUniforms(shaderText string) {
	isUniform := regexp.MustCompile(`uniform\s*(\S*)\s(\S*);`)
	isArray := regexp.MustCompile(`(\w+)\[(\d+)\]`)

	uniformStructs := s.findUniformStructs(shaderText)

	for _, line := range strings.Split(shaderText, "\n") {
		for _, i := range isUniform.FindAllStringSubmatch(line, -1) {
			if len(i) != 3 {
				continue
			}
			if s.resource.UniformExists(i[2]) {
				continue
			}

			name := i[2]
			uType := i[1]

			numItems := 0
			arrayMatch := isArray.FindStringSubmatch(name)
			if len(arrayMatch) != 0 {
				max, err := strconv.Atoi(arrayMatch[2])
				if err != nil {
					fmt.Printf("could not parse '%s' as an integer, %v", arrayMatch[2], err)
				}
				numItems = max
				name = arrayMatch[1]
			}

			if numItems == 0 {
				s.resource.AddUniformName(name)
				s.resource.AdduniformType(uType)
				s.findUniformLocation(uType, i[2], uniformStructs)
				continue
			}

			for k := 0; k < numItems; k++ {
				newName := fmt.Sprintf("%s[%d]", name, k)
				s.resource.AddUniformName(newName)
				s.resource.AdduniformType(uType)
				s.findUniformLocation(uType, newName, uniformStructs)
			}

		}
	}
}

func (s *Shader) findUniformLocation(glType, name string, structs map[string][]glslStruct) {
	structComponents, ok := structs[glType]
	if ok {
		for _, v := range structComponents {
			s.findUniformLocation(v.propertyType, name+"."+v.propertyName, structs)
		}
		return
	}
	t := gl.GetUniformLocation(s.resource.Program, gl.Str(name+"\x00"))
	if t < 0 {
		fmt.Printf("Could not get uniform location for '%s' in shader '%s' (not used?)\n", name, s.filename)
	}
	s.resource.uniforms[name] = t
}

type glslStruct struct {
	propertyName string
	propertyType string
}

func (s Shader) findUniformStructs(shaderText string) map[string][]glslStruct {
	result := make(map[string][]glslStruct)
	isStruct := regexp.MustCompile(`(?s)struct\s*(\w*)\s+{\s([^}]*)};`)
	isStructDefiniton := regexp.MustCompile(`(?s)\s*(\w*)\s(\w*);`)
	for _, structMatch := range isStruct.FindAllStringSubmatch(shaderText, -1) {
		structName := structMatch[1]
		content := structMatch[2]
		var properties []glslStruct
		for _, innerMatch := range isStructDefiniton.FindAllStringSubmatch(content, -1) {
			properties = append(properties, glslStruct{
				propertyName: innerMatch[2],
				propertyType: innerMatch[1],
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
