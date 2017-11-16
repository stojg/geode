package core

import (
	"fmt"
	"github.com/go-gl/gl/v4.1-core/gl"
	"io/ioutil"
	"os"
	"strings"
)

var loadedShaders = make(map[string]*ShaderResource)

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

	// AddAllAttributes(vertexShaderText);

	s.CompileShader()

	//AddAllUniforms(vertexShaderText);
	//AddAllUniforms(fragmentShaderText);

	loadedShaders[fileName] = s.resource

	return s
}

type Shader struct {
	filename string
	resource *ShaderResource
}

func (s *Shader) Bind() {
	gl.UseProgram(s.resource.Program)
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

	//glLogShader(Program, vertex, frag)
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

func NewShaderResource() *ShaderResource {
	s := &ShaderResource{
		Program:  gl.CreateProgram(),
		refCount: 1,
	}

	if s.Program == 0 {
		fmt.Println("Shader creation failed: Could not find valid memory location in constructor")
		os.Exit(1)
	}
	return s
}

type ShaderResource struct {
	Program  uint32
	refCount int
	//uniforms map[string]int
	//uniformNames []string
	//uniformTypes []string
}

func (r *ShaderResource) AddReference() {
	r.refCount++
}

func (r *ShaderResource) Cleanup() {
	//gl.DeleteBuffers(r.Program)
}
