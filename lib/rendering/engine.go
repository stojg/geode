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
		width:      int32(width),
		height:     int32(height),
		samplerMap: samplerMap,
		textures:   make(map[string]components.Texture),
		uniforms:   make(map[string]mgl32.Vec3),

		screenQuad:  NewScreenQuad(),
		nullShader:  NewShader("filter_null"),
		gaussShader: NewShader("filter_gauss"),

		ambientShader: NewShader("forward_ambient"),

		hdrTexture:    framebuffer.NewTexture(0, width, height, gl.RGBA32F, gl.RGBA, gl.FLOAT, gl.LINEAR, false),
		toneMapShader: NewShader("filter_tonemap"),

		shadowMapTemp:  framebuffer.NewTexture(0, 512*2, 512*2, gl.RG32F, gl.RGB, gl.FLOAT, gl.LINEAR, true),
		fullscreenTemp: framebuffer.NewTexture(0, width, height, gl.RGBA32F, gl.RGBA, gl.FLOAT, gl.LINEAR, false),
	}

	e.SetTexture("x_shadowMap", e.shadowMapTemp)

	return e
}

type Engine struct {
	width, height int32
	mainCamera    *components.Camera
	lights        []components.Light
	activeLight   components.Light

	samplerMap map[string]uint32
	textures   map[string]components.Texture
	uniforms   map[string]mgl32.Vec3

	screenQuad    *ScreenQuad
	nullShader    *Shader
	gaussShader   *Shader
	ambientShader *Shader
	toneMapShader *Shader

	hdrTexture *framebuffer.Texture

	shadowMapTemp  *framebuffer.Texture
	fullscreenTemp *framebuffer.Texture
}

func (e *Engine) Render(object components.Renderable) {
	if e.mainCamera == nil {
		panic("mainCamera not found, the game cannot render")
	}
	checkForError("renderer.Engine.Render [start]")

	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	// shadow map

	gl.Enable(gl.DEPTH_TEST)
	for _, l := range e.lights {
		caster, ok := l.(components.ShadowCaster)
		if !ok {
			continue
		}
		e.activeLight = l
		e.SetTexture("x_shadowMap", caster.ShadowTexture())
		caster.ShadowTexture().BindAsRenderTarget()

		gl.CullFace(gl.FRONT)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		object.RenderAll(caster.ShadowShader(), e)
		gl.CullFace(gl.BACK)
		//
		e.blurShadowMap(caster.ShadowTexture(), 1)
		//
		//	// debug
		//	//gl.Viewport(0, 0, e.width, e.height)
		//	//gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
		//	//gl.Disable(gl.DEPTH_TEST)
		//	//caster.BindShadow()
		//	//e.screenShader.Bind()
		//	//gl.Clear(gl.COLOR_BUFFER_BIT)
		//	//e.screenQuad.Draw()
		//	//return
	}

	e.hdrTexture.BindAsRenderTarget()
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	e.hdrTexture.SetViewPort()

	// ambient pass
	object.RenderAll(e.ambientShader, e)

	// light pass
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.ONE, gl.ONE)
	gl.DepthMask(false)
	gl.DepthFunc(gl.EQUAL)

	for _, l := range e.lights {
		e.activeLight = l
		//l.Shader().Bind()
		if caster, ok := l.(components.ShadowCaster); ok {
			e.SetTexture("x_shadowMap", caster.ShadowTexture())
		}
		object.RenderAll(l.Shader(), e)
	}
	gl.DepthFunc(gl.LESS)
	gl.DepthMask(true)
	gl.Disable(gl.BLEND)

	e.applyFilter(e.toneMapShader, e.hdrTexture, nil)

	// move to default framebuffer buffer
	//gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
	//gl.Viewport(0, 0, e.width, e.height)
	// disable depth test so screen-space quad isn't discarded due to depth test
	//gl.Disable(gl.DEPTH_TEST)
	//e.toneMapShader.Bind()
	//e.hdrTexture.Bind(e.GetSamplerSlot("x_filterTexture"))
	//gl.Clear(gl.COLOR_BUFFER_BIT)
	//e.screenQuad.Draw()

	checkForError("renderer.Engine.Render [end]")
	//os.Exit(0)
}

func (e *Engine) GetActiveLight() components.Light {
	return e.activeLight
}

func (e *Engine) AddLight(l components.Light) {
	e.lights = append(e.lights, l)
}

func (e *Engine) blurShadowMap(shadowMap components.Texture, blurAmount float32) {

	e.SetVector3f("x_blurScale", mgl32.Vec3{1 / float32(shadowMap.Width()) * blurAmount, 0, 0})
	e.applyFilter(e.gaussShader, shadowMap, e.shadowMapTemp)
	e.SetVector3f("x_blurScale", mgl32.Vec3{0, 1 / float32(shadowMap.Height()) * blurAmount, 0})
	e.applyFilter(e.gaussShader, e.shadowMapTemp, shadowMap)
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
	// m_altCamera.SetProjection(Matrix4f().InitIdentity());
	// m_altCamera.GetTransform()->SetPos(Vector3f(0,0,0));
	// m_altCamera.GetTransform()->SetRot(Quaternion(Vector3f(0,1,0),ToRadians(180.0f)));

	gl.Clear(gl.DEPTH_BUFFER_BIT)

	//id := e.GetSamplerSlot("x_filterTexture")
	//in.Bind(id)
	filter.Bind()
	filter.UpdateUniforms(nil, nil, e)

	e.screenQuad.Draw()
	//e.SetTexture("x_filterTexture", nil)
}

/*
	inline void SetVector3f(const std::string& name, const Vector3f& value) { m_vector3fMap[name] = value; }
	inline void SetFloat(const std::string& name, float value)              { m_floatMap[name] = value; }
	inline void SetTexture(const std::string& name, const Texture& value)   { m_textureMap[name] = value; }

	const Vector3f& GetVector3f(const std::string& name) const;
	float GetFloat(const std::string& name)              const;
	const Texture& GetTexture(const std::string& name)   const;
*/

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

func (e *Engine) SetVector3f(name string, v mgl32.Vec3) {
	e.uniforms[name] = v
}

func (e *Engine) GetVector3f(name string) mgl32.Vec3 {
	v, ok := e.uniforms[name]
	if !ok {
		panic(fmt.Sprintf("no value found for uniform '%s'", name))
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
