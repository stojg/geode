package rendering

import (
	"fmt"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/stojg/graphics/lib/components"
	"github.com/stojg/graphics/lib/rendering/framebuffer"
)

const maxShadowMaps = 3

func NewEngine(width, height int) *Engine {

	gl.ClearColor(0.00, 0.00, 0.00, 1)

	gl.FrontFace(gl.CCW)
	gl.CullFace(gl.BACK)
	gl.Enable(gl.CULL_FACE)
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)

	gl.Disable(gl.MULTISAMPLE)
	gl.Disable(gl.FRAMEBUFFER_SRGB)

	samplerMap := make(map[string]uint32)
	samplerMap["diffuse"] = 0
	samplerMap["x_shadowMap"] = 9
	samplerMap["x_filterTexture"] = 10

	e := &Engine{
		width:      int32(width),
		height:     int32(height),
		samplerMap: samplerMap,
		textures:   make(map[string]components.Texture),
		uniforms3f: make(map[string]mgl32.Vec3),
		uniformsi:  make(map[string]int32),

		screenQuad: NewScreenQuad(),

		nullShader:    NewShader("filter_null"),
		fxaaShader:    NewShader("filter_fxaa"),
		gaussShader:   NewShader("filter_gauss"),
		ambientShader: NewShader("forward_ambient"),
		shadowShader:  NewShader("shadow_vsm"),

		screenTexture: framebuffer.NewTexture(gl.COLOR_ATTACHMENT0, width, height, gl.RGBA32F, gl.RGBA, gl.FLOAT, gl.NEAREST, false),
		toneMapShader: NewShader("filter_tonemap"),

		fullScreenTemp: framebuffer.NewTexture(gl.COLOR_ATTACHMENT0, width, height, gl.RGBA32F, gl.RGBA, gl.FLOAT, gl.NEAREST, false),

		capabilities: make(map[string]bool),
	}

	shadowW, shadowH := 1024, 1024
	checkForError("rendering.NewEngine end")
	for i := 0; i < maxShadowMaps; i++ {
		e.shadowTextures[i] = framebuffer.NewTexture(gl.COLOR_ATTACHMENT0, shadowW, shadowH, gl.RG32F, gl.RGB, gl.FLOAT, gl.LINEAR, true)
	}
	e.tempShadowTexture = framebuffer.NewTexture(gl.COLOR_ATTACHMENT0, shadowW, shadowH, gl.RG32F, gl.RGB, gl.FLOAT, gl.LINEAR, true)
	return e
}

type Engine struct {
	width, height int32
	mainCamera    *components.Camera
	lights        []components.Light
	activeLight   components.Light

	samplerMap map[string]uint32
	textures   map[string]components.Texture
	uniforms3f map[string]mgl32.Vec3
	uniformsi  map[string]int32

	screenQuad    *ScreenQuad
	nullShader    *Shader
	gaussShader   *Shader
	ambientShader *Shader
	toneMapShader *Shader
	shadowShader  *Shader
	fxaaShader    *Shader

	screenTexture *framebuffer.Texture

	shadowTextures    [maxShadowMaps]components.Texture
	tempShadowTexture components.Texture

	fullScreenTemp *framebuffer.Texture

	capabilities map[string]bool
}

func (e *Engine) Enable(cap string) {
	e.capabilities[cap] = true
}

func (e *Engine) Disable(cap string) {
	e.capabilities[cap] = false
}
func (e *Engine) Render(object components.Renderable) {
	if e.mainCamera == nil {
		panic("mainCamera not found, the game cannot render")
	}
	checkForError("renderer.Engine.Render [start]")

	// shadow map
	gl.Enable(gl.DEPTH_TEST)
	gl.ClearColor(1, 1, 1, 1)
	for i, l := range e.lights {
		e.activeLight = l
		if !l.ShadowCaster() {
			continue
		}
		e.shadowTextures[i].BindAsRenderTarget()
		gl.Clear(gl.DEPTH_BUFFER_BIT | gl.COLOR_BUFFER_BIT)
		object.RenderAll(e.shadowShader, e)

		e.blurShadowMap(e.shadowTextures[i], 1)
	}

	//gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
	e.screenTexture.BindAsRenderTarget()
	e.screenTexture.SetViewPort()
	gl.ClearColor(0.541, 0.616, 0.671, 1)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	// ambient pass
	object.RenderAll(e.ambientShader, e)

	// light pass
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.ONE, gl.ONE)
	gl.DepthMask(false)
	gl.DepthFunc(gl.EQUAL)

	for i, l := range e.lights {
		e.activeLight = l
		if l.ShadowCaster() {
			e.SetTexture("x_shadowMap", e.shadowTextures[i])
		}
		object.RenderAll(l.Shader(), e)
	}
	gl.DepthFunc(gl.LESS)
	gl.DepthMask(true)
	gl.Disable(gl.BLEND)

	e.applyFilter(e.toneMapShader, e.screenTexture, e.fullScreenTemp)
	e.applyFilter(e.fxaaShader, e.fullScreenTemp, nil)

	checkForError("renderer.Engine.Render [end]")
}

func (e *Engine) GetActiveLight() components.Light {
	return e.activeLight
}

func (e *Engine) AddLight(l components.Light) {
	e.lights = append(e.lights, l)
}

func (e *Engine) blurShadowMap(shadowMap components.Texture, blurAmount float32) {
	e.SetVector3f("x_blurScale", mgl32.Vec3{1 / float32(shadowMap.Width()) * blurAmount, 0, 0})
	e.applyFilter(e.gaussShader, shadowMap, e.tempShadowTexture)
	e.SetVector3f("x_blurScale", mgl32.Vec3{0, 1 / float32(shadowMap.Height()) * blurAmount, 0})
	e.applyFilter(e.gaussShader, e.tempShadowTexture, shadowMap)
}

func (e *Engine) applyFilter(filter *Shader, in, out components.Texture) {
	if in == out {
		panic("Argh, can apply filter where source and destination is the same")
	}

	if out == nil {
		gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
	} else {
		out.BindAsRenderTarget()
	}
	e.SetTexture("x_filterTexture", in)
	gl.Clear(gl.DEPTH_BUFFER_BIT)
	filter.Bind()
	filter.UpdateUniforms(nil, nil, e)
	e.screenQuad.Draw()
}

func (e *Engine) SetTexture(name string, texture components.Texture) {
	e.textures[name] = texture
}

func (e *Engine) GetTexture(name string) components.Texture {
	v, ok := e.textures[name]
	if !ok {
		fmt.Printf("Could not find texture '%s'\n", name)
		panic("")
	}
	return v
}

func (e *Engine) SetInteger(name string, v int32) {
	e.uniformsi[name] = v
}

func (e *Engine) GetInteger(name string) int32 {
	v, ok := e.uniformsi[name]
	if !ok {
		panic(fmt.Sprintf("GetInteger, no value found for uniform '%s'", name))
	}
	return v
}

func (e *Engine) SetVector3f(name string, v mgl32.Vec3) {
	e.uniforms3f[name] = v
}

func (e *Engine) GetVector3f(name string) mgl32.Vec3 {
	v, ok := e.uniforms3f[name]
	if !ok {
		panic(fmt.Sprintf("GetVector3f, no value found for uniform '%s'", name))
	}
	return v
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
