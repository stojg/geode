package rendering

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/stojg/graphics/lib/components"
)

type GameObject interface {
	RenderAll(*Shader, components.UniformUpdater)
}

func NewEngine() *Engine {

	gl.ClearColor(0.05, 0.06, 0.07, 0)

	gl.FrontFace(gl.CCW)
	gl.CullFace(gl.BACK)
	gl.Enable(gl.CULL_FACE)
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)

	gl.Enable(gl.MULTISAMPLE)
	gl.Enable(gl.FRAMEBUFFER_SRGB)

	s := NewShader("simple")

	return &Engine{
		shader: s,
	}
}

type Engine struct {
	shader *Shader
}

func (e *Engine) Render(object GameObject) {
	CheckForError("renderer.Engine.Render [start]")
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	object.RenderAll(e.shader, e)
	CheckForError("renderer.Engine.Render [end]")
}
