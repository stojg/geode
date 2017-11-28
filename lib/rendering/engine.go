package rendering

import (
	"fmt"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/stojg/graphics/lib/components"
	"github.com/stojg/graphics/lib/rendering/framebuffer"
)

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
		width:         int32(width),
		height:        int32(height),
		samplerMap:    samplerMap,
		textures:      make(map[string]components.Texture),
		uniforms3f:    make(map[string]mgl32.Vec3),
		uniformsI:     make(map[string]int32),
		uniformsFloat: make(map[string]float32),

		screenQuad: NewScreenQuad(),

		nullShader:    NewShader("filter_null"),
		fxaaShader:    NewShader("filter_fxaa"),
		gaussShader:   NewShader("filter_gauss"),
		ambientShader: NewShader("forward_ambient"),
		shadowShader:  NewShader("shadow_vsm"),

		debugShadowShader: NewShader("debug_shadow"),

		screenTexture: framebuffer.NewTexture(gl.COLOR_ATTACHMENT0, width, height, gl.RGB16F, gl.RGB, gl.FLOAT, gl.NEAREST, false),
		toneMapShader: NewShader("filter_tonemap"),

		fullScreenTemp: framebuffer.NewTexture(gl.COLOR_ATTACHMENT0, width, height, gl.RGB, gl.RGB, gl.UNSIGNED_BYTE, gl.NEAREST, false),

		capabilities: make(map[string]bool),
	}

	e.shadowTextures = make([]components.Texture, 11)
	e.tempShadowTextures = make([]components.Texture, 11)
	e.shadowTextures[0] = framebuffer.NewTexture(gl.COLOR_ATTACHMENT0, 1, 1, gl.RG32F, gl.RG, gl.FLOAT, gl.LINEAR_MIPMAP_LINEAR, true)
	e.tempShadowTextures[0] = framebuffer.NewTexture(gl.COLOR_ATTACHMENT0, 1, 1, gl.RG32F, gl.RG, gl.FLOAT, gl.LINEAR_MIPMAP_LINEAR, true)
	for i := uint(0); i < 10; i++ {
		size := 2 << i // power of two, 2, 4, 8, 16 and so on
		e.shadowTextures[i+1] = framebuffer.NewTexture(gl.COLOR_ATTACHMENT0, size, size, gl.RG32F, gl.RG, gl.FLOAT, gl.LINEAR_MIPMAP_LINEAR, true)
		e.tempShadowTextures[i+1] = framebuffer.NewTexture(gl.COLOR_ATTACHMENT0, size, size, gl.RG32F, gl.RG, gl.FLOAT, gl.LINEAR_MIPMAP_LINEAR, true)
	}

	// set defaults
	e.SetFloat("x_varianceMin", 0.0)
	e.SetFloat("x_lightBleedReductionAmount", 0.0)
	e.SetTexture("x_shadowMap", e.shadowTextures[0])

	checkForError("rendering.NewEngine end")
	return e
}

type Engine struct {
	width, height int32
	mainCamera    components.Viewable
	lights        []components.Light
	activeLight   components.Light

	samplerMap    map[string]uint32
	textures      map[string]components.Texture
	uniforms3f    map[string]mgl32.Vec3
	uniformsI     map[string]int32
	uniformsFloat map[string]float32

	screenQuad    *ScreenQuad
	nullShader    *Shader
	gaussShader   *Shader
	ambientShader *Shader
	toneMapShader *Shader
	shadowShader  *Shader
	fxaaShader    *Shader

	debugShadowShader *Shader

	screenTexture *framebuffer.Texture

	shadowTextures     []components.Texture
	tempShadowTextures []components.Texture

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
	gl.Enable(gl.DEPTH_TEST)

	// shadow map
	for _, l := range e.lights {
		e.activeLight = l
		if !l.ShadowCaster() {
			continue
		}
		info := l.ShadowInfo()

		e.shadowTextures[info.SizeAsPowerOfTwo()].BindAsRenderTarget()
		e.shadowTextures[info.SizeAsPowerOfTwo()].SetViewPort()
		gl.Clear(gl.DEPTH_BUFFER_BIT | gl.COLOR_BUFFER_BIT)

		if info.FlipFaces() {
			gl.CullFace(gl.FRONT)
		}

		object.RenderAll(e.shadowShader, e)

		if info.FlipFaces() {
			gl.CullFace(gl.BACK)
		}

		gl.GenerateMipmap(gl.TEXTURE_2D)

		//gl.Disable(gl.DEPTH_TEST)
		//gl.Viewport(0, 0, e.width, e.height)
		//e.applyFilter(e.debugShadowShader, e.shadowTextures[i], nil)
		//return

		e.blurShadowMap(info.SizeAsPowerOfTwo(), 1)
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

	for _, l := range e.lights {
		e.activeLight = l
		if l.ShadowCaster() {
			info := l.ShadowInfo()
			e.SetFloat("x_varianceMin", info.ShadowVarianceMin())
			e.SetFloat("x_lightBleedReductionAmount", info.LightBleedReduction())
			e.SetTexture("x_shadowMap", e.shadowTextures[info.SizeAsPowerOfTwo()])
		}
		object.RenderAll(l.Shader(), e)
	}
	gl.DepthFunc(gl.LESS)
	gl.DepthMask(true)
	gl.Disable(gl.BLEND)

	gl.Disable(gl.DEPTH_TEST)
	e.applyFilter(e.toneMapShader, e.screenTexture, e.fullScreenTemp)
	e.applyFilter(e.fxaaShader, e.fullScreenTemp, nil)
	gl.Enable(gl.DEPTH_TEST)

	checkForError("renderer.Engine.Render [end]")
}

func (e *Engine) ActiveLight() components.Light {
	return e.activeLight
}

func (e *Engine) AddLight(l components.Light) {
	e.lights = append(e.lights, l)
}

func (e *Engine) blurShadowMap(sizeAsPowerOfTwo int, blurAmount float32) {
	var size = 2 << uint(sizeAsPowerOfTwo)
	gl.Disable(gl.DEPTH_TEST)
	e.SetVector3f("x_blurScale", mgl32.Vec3{1 / float32(size) * blurAmount, 0, 0})
	e.applyFilter(e.gaussShader, e.shadowTextures[sizeAsPowerOfTwo], e.tempShadowTextures[sizeAsPowerOfTwo])
	e.SetVector3f("x_blurScale", mgl32.Vec3{0, 1 / float32(size) * blurAmount, 0})
	e.applyFilter(e.gaussShader, e.tempShadowTextures[sizeAsPowerOfTwo], e.shadowTextures[sizeAsPowerOfTwo])
	gl.GenerateMipmap(gl.TEXTURE_2D)
	gl.Enable(gl.DEPTH_TEST)
}

func (e *Engine) applyFilter(filter *Shader, in, out components.Texture) {
	if in == out {
		panic("Argh, can't apply filter where source and destination is the same")
	}

	if out == nil {
		gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
	} else {
		out.BindAsRenderTarget()
	}
	e.SetTexture("x_filterTexture", in)
	filter.Bind()
	filter.UpdateUniforms(nil, nil, e)
	e.screenQuad.Draw()
}

func (e *Engine) SetTexture(name string, texture components.Texture) {
	e.textures[name] = texture
}

func (e *Engine) Texture(name string) components.Texture {
	v, ok := e.textures[name]
	if !ok {
		panic(fmt.Sprintf("Texture, Could not find texture '%s'\n", name))
	}
	return v
}

func (e *Engine) SetInteger(name string, v int32) {
	e.uniformsI[name] = v
}

func (e *Engine) Integer(name string) int32 {
	v, ok := e.uniformsI[name]
	if !ok {
		panic(fmt.Sprintf("Integer, no value found for uniform '%s'", name))
	}
	return v
}

func (e *Engine) SetFloat(name string, v float32) {
	e.uniformsFloat[name] = v
}

func (e *Engine) Float(name string) float32 {
	v, ok := e.uniformsFloat[name]
	if !ok {
		panic(fmt.Sprintf("Float, no value found for uniform '%s'", name))
	}
	return v
}

func (e *Engine) SetVector3f(name string, v mgl32.Vec3) {
	e.uniforms3f[name] = v
}

func (e *Engine) Vector3f(name string) mgl32.Vec3 {
	v, ok := e.uniforms3f[name]
	if !ok {
		panic(fmt.Sprintf("Vector3f, no value found for uniform '%s'", name))
	}
	return v
}

func (e *Engine) AddCamera(c components.Viewable) {
	e.mainCamera = c
}

func (e *Engine) MainCamera() components.Viewable {
	return e.mainCamera
}

func (e *Engine) SamplerSlot(samplerName string) uint32 {
	slot, exists := e.samplerMap[samplerName]
	if !exists {
		fmt.Printf("rendering.Engine tried finding texture slot for %s, failed\n", samplerName)
	}
	return slot
}
