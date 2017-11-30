package shader

import (
	"testing"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

func TestShader_NewShader(t *testing.T) {

	defer setupOpenGL(t)()

	shaderFolder = "./testdata/"

	s := NewShader("test")

	if s.filename != "test" {
		t.Errorf("Shader.filename not set correctly, wanted %s, got %s", "forward_point", s.filename)
	}

	expectedUniformNames := []string{
		"projection", "view", "model", "light", "diffuse",
	}
	expectedUniformTypes := []string{
		"mat4", "mat4", "mat4", "TestStructB", "sampler2D",
	}

	expectedUniforms := []string{
		"light.inner.color", "light.position", "diffuse", "projection", "view", "model",
	}

	if len(s.resource.uniformNames) != len(expectedUniformNames) {
		t.Errorf("expected %d uniform names, got %d", len(expectedUniformNames), len(s.resource.uniformNames))
		return
	}

	for i, expected := range expectedUniformNames {
		if expected != s.resource.uniformNames[i] {
			t.Errorf("expected uniform name %s, got %s", expected, s.resource.uniformNames[i])
		}
	}

	if len(s.resource.uniformTypes) != len(expectedUniformTypes) {
		t.Errorf("expected %d uniform types, got %d", len(expectedUniformNames), len(s.resource.uniformTypes))
		return
	}

	for i, expected := range expectedUniformTypes {
		if expected != s.resource.uniformTypes[i] {
			t.Errorf("expected uniform type %s, got %s", expected, s.resource.uniformTypes[i])
		}
	}

	if len(s.resource.uniforms) != len(expectedUniforms) {
		t.Errorf("expected %d uniforms, got %d", len(expectedUniforms), len(s.resource.uniforms))
		return
	}

	for _, expectedUniform := range expectedUniforms {
		if _, ok := s.resource.uniforms[expectedUniform]; !ok {
			t.Errorf("expected uniform %s, but could not find it", expectedUniform)
		}
	}
}

func setupOpenGL(t *testing.T) func() {
	//runtime.LockOSThread()
	if err := glfw.Init(); err != nil {
		t.Errorf("failed to initialize glfw: %s", err)
		return func() {}
	}

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	window, err := glfw.CreateWindow(1, 1, "test", nil, nil)
	if err != nil {
		t.Error(err)
		return func() {}
	}
	window.MakeContextCurrent()

	if err := gl.Init(); err != nil {
		t.Error(err)
		return func() {}
	}

	return func() {
		window.Destroy()
		glfw.Terminate()
	}
}
