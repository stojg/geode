package shader

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/stojg/geode/lib/components"
	"github.com/stojg/geode/lib/debug"
	"github.com/stojg/geode/lib/physics"
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

	// @todo, the information about which index this UBO is in the renderstate
	uboID := gl.GetUniformBlockIndex(s.resource.Program, gl.Str("Matrices\x00"))
	if uboID != math.MaxUint32 {
		gl.UniformBlockBinding(s.resource.Program, uboID, 0)
	}
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
		debug.ShaderSwitch()
	}
}

func (s *Shader) Unbind() {
	gl.UseProgram(0)
}

func (s *Shader) UpdateTransform(transform *physics.Transform, engine components.RenderState) {
	for _, name := range s.resource.uniformNames {
		switch name {
		case "LightVP":
			s.UpdateUniform(name, engine.ActiveLight().ViewProjection())
		case "model":
			s.UpdateUniform(name, transform.Transformation())
		}
	}
}

func (s *Shader) UpdateUniforms(material components.Material, state components.RenderState) {

	for i, name := range s.resource.uniformNames {
		uniformType := s.resource.uniformTypes[i]

		// x_ means that the uniform is a state uniform and should be fetch from the "global" state
		if strings.Index(name, "x_") == 0 {
			s.updateUniformFromState(uniformType, state, name)
			continue
		}

		if uniformType == "sampler2D" || uniformType == "samplerCube" {
			samplerSlot := state.SamplerSlot(name)
			material.Texture(name).Activate(samplerSlot)
			s.UpdateUniform(name, int32(samplerSlot))
		}

		name, index := getArray(name)
		switch name {
		case "lights":
			if len(state.Lights()) > index {
				light := state.Lights()[index]
				s.UpdateUniform(fmt.Sprintf("%s[%d].position", name, index), light.Position())
				s.UpdateUniform(fmt.Sprintf("%s[%d].color", name, index), light.Color())
				s.UpdateUniform(fmt.Sprintf("%s[%d].constant", name, index), light.Constant())
				s.UpdateUniform(fmt.Sprintf("%s[%d].linear", name, index), light.Linear())
				s.UpdateUniform(fmt.Sprintf("%s[%d].quadratic", name, index), light.Exponent())
				s.UpdateUniform(fmt.Sprintf("%s[%d].distance", name, index), light.MaxDistance())
				s.UpdateUniform(fmt.Sprintf("%s[%d].direction", name, index), light.Direction())
				s.UpdateUniform(fmt.Sprintf("%s[%d].cutoff", name, index), light.Cutoff())
			}
		case "numLights":
			s.UpdateUniform(name, int32(len(state.Lights())))
		}
	}
}

func (s *Shader) UpdateUniform(uniformName string, value interface{}) {
	loc, ok := s.resource.uniforms[uniformName]
	if !ok {
		panic(fmt.Sprintf("no shader location found for uniform: '%s' in shader '%s'", uniformName, s.filename))
	}
	debug.AddUniformSet()
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

func (s *Shader) updateUniformFromState(uniformType string, engine components.RenderState, name string) {
	if uniformType == "sampler2D" || uniformType == "samplerCube" {
		samplerSlot := engine.SamplerSlot(name)
		engine.Texture(name).Activate(samplerSlot)
		gl.Uniform1i(s.resource.uniforms[name], int32(samplerSlot))
		debug.AddUniformSet()
		return
	}

	switch uniformType {
	case "bool":
		s.UpdateUniform(name, engine.Integer(name))
	case "vec3":
		s.UpdateUniform(name, engine.Vector3f(name))
	case "float":
		s.UpdateUniform(name, engine.Float(name))
	case "int":
		s.UpdateUniform(name, engine.Integer(name))
	default:
		panic(fmt.Sprintf("Shader.UpdateUniforms() don't know how to set uniformType '%s' with name '%s'", uniformType, name))
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

	var re = regexp.MustCompile(`^#include\s"([^"]*)"`)

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

func (s *Shader) setUniformBaseLight(uniformName string, baseLight components.Light) {
	s.UpdateUniform(uniformName+".color", baseLight.Color())
	s.UpdateUniform(uniformName+".maxDistance", baseLight.MaxDistance())
}

func (s *Shader) addVertexShader(shader string) uint32 {
	return s.createProgram(shader, gl.VERTEX_SHADER)
}

func (s *Shader) addFragmentShader(shader string) uint32 {
	return s.createProgram(shader, gl.FRAGMENT_SHADER)
}

var isArray = regexp.MustCompile(`(\w+)\[(\d+)\]`)

func (s *Shader) findAllUniforms(shaderText string) {
	isUniform := regexp.MustCompile(`^uniform\s+(\S*)\s+(\S*);`)

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

			name, numItems := getArray(name)

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

func getArray(name string) (string, int) {
	number := 0
	arrayMatch := isArray.FindStringSubmatch(name)
	if len(arrayMatch) != 0 {
		max, err := strconv.Atoi(arrayMatch[2])
		if err != nil {
			fmt.Printf("could not parse '%s' as an integer, %v", arrayMatch[2], err)
		}
		return arrayMatch[1], max
	}
	return name, number
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
		fmt.Printf("%d: %s\n", i+1, line)
	}
}
