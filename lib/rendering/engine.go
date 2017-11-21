package rendering

import (
	"fmt"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/stojg/graphics/lib/components"
	"github.com/stojg/graphics/lib/rendering/framebuffer"
)

type GameObject interface {
	RenderAll(shader components.Shader, engine components.RenderingEngine)
}

func NewEngine(width, height int) *Engine {

	gl.ClearColor(0.04, 0.04, 0.04, 0)

	gl.FrontFace(gl.CCW)
	gl.CullFace(gl.BACK)
	gl.Enable(gl.CULL_FACE)
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)

	gl.Enable(gl.MULTISAMPLE)
	gl.Disable(gl.FRAMEBUFFER_SRGB)

	samplerMap := make(map[string]uint32)
	samplerMap["diffuse"] = 0
	samplerMap["shadowMap"] = 1

	return &Engine{
		width:      width,
		height:     height,
		samplerMap: samplerMap,

		screenQuad:   NewScreenQuad(),
		screenShader: NewShader("screen_shader"),

		ambientShader: NewShader("forward_ambient"),

		hdrBuffer: framebuffer.NewHDR(int32(width), int32(height)),
		hdrShader: NewShader("screen_hdr"),

		shadowBuffer: framebuffer.NewShadow(1024, 1024),
		shadowShader: NewShader("shadow"),
	}
}

type Engine struct {
	width, height int
	mainCamera    *components.Camera
	lights        []components.Light
	activeLight   components.Light

	samplerMap map[string]uint32

	screenQuad   *ScreenQuad
	screenShader *Shader

	ambientShader *Shader

	hdrBuffer *framebuffer.FBO
	hdrShader *Shader

	shadowBuffer *framebuffer.FBO
	shadowShader *Shader
}

func (e *Engine) Render(object GameObject) {
	if e.mainCamera == nil {
		panic("mainCamera not found, the game cannot render")
	}
	CheckForError("renderer.Engine.Render [start]")

	// shadow map
	{
		gl.Viewport(0, 0, 1024, 1024)
		e.shadowBuffer.Bind()
		gl.Enable(gl.DEPTH_TEST)
		gl.Clear(gl.DEPTH_BUFFER_BIT)
		object.RenderAll(e.shadowShader, e)

		//// debug
		//gl.Viewport(0, 0, int32(e.width), int32(e.height))
		//gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
		//
		//gl.Disable(gl.DEPTH_TEST)
		//e.shadowBuffer.BindTexture()
		//
		//e.screenShader.Bind()
		//gl.Clear(gl.COLOR_BUFFER_BIT)
		//e.screenQuad.Draw()
		//return
	}

	// ambient pass
	gl.Viewport(0, 0, int32(e.width), int32(e.height))
	e.hdrBuffer.Bind()
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.Enable(gl.DEPTH_TEST)

	object.RenderAll(e.ambientShader, e)

	// light pass
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.ONE, gl.ONE)
	gl.DepthMask(false)
	gl.DepthFunc(gl.EQUAL)
	for _, l := range e.lights {
		e.activeLight = l

		l.Shader().Bind()
		//mat.Texture(name).Bind(samplerSlot)

		gl.ActiveTexture(gl.TEXTURE0 + uint32(1))
		l.Shader().SetUniformi("shadowMap", int32(1))
		gl.BindTexture(gl.TEXTURE_2D, e.shadowBuffer.Texture().ID())

		object.RenderAll(l.Shader(), e)
	}
	gl.DepthFunc(gl.LESS)
	gl.DepthMask(true)
	gl.Disable(gl.BLEND)

	// move to default framebuffer buffer
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)

	// disable depth test so screen-space quad isn't discarded due to depth test
	gl.Disable(gl.DEPTH_TEST)
	e.hdrBuffer.BindTexture()
	e.hdrShader.Bind()
	gl.Clear(gl.COLOR_BUFFER_BIT)
	e.screenQuad.Draw()

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

func (e *Engine) GetSamplerSlot(samplerName string) uint32 {
	slot, exists := e.samplerMap[samplerName]
	if !exists {
		fmt.Printf("rendering.Engine tried finding texture slot for %s, failed\n", samplerName)
	}
	return slot
}
