package rendering

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/stojg/graphics/lib/components"
)

type GameObject interface {
	RenderAll(shader components.Shader, engine components.RenderingEngine)
}

func NewEngine() *Engine {

	gl.ClearColor(0.04, 0.04, 0.04, 0)

	gl.FrontFace(gl.CCW)
	gl.CullFace(gl.BACK)
	gl.Enable(gl.CULL_FACE)
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)

	gl.Enable(gl.MULTISAMPLE)
	gl.Enable(gl.FRAMEBUFFER_SRGB)

	s := NewShader("forward_ambient")

	return &Engine{
		ambientShader: s,
	}
}

type Engine struct {
	ambientShader *Shader
	mainCamera    *components.Camera
	lights        []components.Light
	activeLight   components.Light
}

func (e *Engine) Render(object GameObject) {
	if e.mainCamera == nil {
		panic("mainCamera not found, the game cannot render")
	}
	CheckForError("renderer.Engine.Render [start]")

	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	object.RenderAll(e.ambientShader, e)

	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.ONE, gl.ONE)
	gl.DepthMask(false)
	gl.DepthFunc(gl.EQUAL)

	for _, l := range e.lights {
		e.activeLight = l
		object.RenderAll(l.Shader(), e)
	}

	gl.DepthFunc(gl.LESS)
	gl.DepthMask(true)
	gl.Disable(gl.BLEND)

	CheckForError("renderer.Engine.Render [end]")
}

func (e *Engine) GetActiveLight() components.Light {
	return e.activeLight
}

func (e *Engine) AddLight(l components.Light) {
	e.lights = append(e.lights, l)
}

func (e *Engine) AddCamera(c *components.Camera) {
	e.mainCamera = c
}

func (e *Engine) GetMainCamera() *components.Camera {
	return e.mainCamera
}
