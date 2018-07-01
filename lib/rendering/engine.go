package rendering

import (
	"fmt"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/stojg/graphics/lib/components"
	"github.com/stojg/graphics/lib/debug"
	"github.com/stojg/graphics/lib/rendering/debugger"
	"github.com/stojg/graphics/lib/rendering/framebuffer"
	"github.com/stojg/graphics/lib/rendering/primitives"
	"github.com/stojg/graphics/lib/rendering/shader"
	"github.com/stojg/graphics/lib/rendering/shadow"
	"github.com/stojg/graphics/lib/rendering/standard"
	"github.com/stojg/graphics/lib/rendering/technique"
	"github.com/stojg/graphics/lib/rendering/terrain"
)

func NewEngine(width, height int) *Engine {

	var nrAttributes int32
	gl.GetIntegerv(gl.MAX_VERTEX_ATTRIBS, &nrAttributes)
	fmt.Printf("maximum nr of vertex attributes supported: %d\n", nrAttributes)

	gl.ClearColor(0.01, 0.01, 0.01, 1)

	gl.FrontFace(gl.CCW)
	gl.CullFace(gl.BACK)
	gl.Enable(gl.CULL_FACE)
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)

	gl.Enable(gl.TEXTURE_CUBE_MAP_SEAMLESS)
	gl.Disable(gl.FRAMEBUFFER_SRGB)

	samplerMap := make(map[string]uint32)
	samplerMap["albedo"] = 0
	samplerMap["metallic"] = 1
	samplerMap["roughness"] = 2
	samplerMap["normal"] = 3
	samplerMap["x_shadowMap"] = 9
	samplerMap["x_filterTexture"] = 10

	e := &Engine{
		width:  int32(width),
		height: int32(height),

		samplerMap: samplerMap,

		textures:      make(map[string]components.Texture),
		uniforms3f:    make(map[string]mgl32.Vec3),
		uniformsI:     make(map[string]int32),
		uniformsFloat: make(map[string]float32),

		nullShader:    shader.NewShader("filter_null"),
		overlayShader: shader.NewShader("filter_overlay"),

		multiSampledTexture: framebuffer.NewMultiSampledTexture(gl.COLOR_ATTACHMENT0, width, height, gl.RGBA16F, gl.RGBA, gl.FLOAT, gl.LINEAR, false),
		offScreenTexture:    framebuffer.NewTexture(gl.COLOR_ATTACHMENT0, width, height, gl.RGBA16F, gl.RGBA, gl.FLOAT, gl.LINEAR, false),
		toneMapShader:       shader.NewShader("filter_tonemap"),

		fullScreenTemp: framebuffer.NewTexture(gl.COLOR_ATTACHMENT0, width, height, gl.RGB, gl.RGB, gl.UNSIGNED_BYTE, gl.NEAREST, false),
	}
	e.standardRenderer = standard.NewRenderer(e)
	e.shadowMap = shadow.NewRenderer(e)
	e.terrainRenderer = terrain.NewRenderer(e)

	e.skybox = technique.NewSkyBox("res/textures/sky0016.hdr", e)

	debugger.New(width, height)

	e.SetInteger("x_enable_env_map", 1)
	e.SetInteger("x_enable_skybox", 1)

	debug.CheckForError("rendering.NewEngine end")
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

	nullShader    *shader.Shader
	gaussShader   *shader.Shader
	toneMapShader *shader.Shader
	overlayShader *shader.Shader

	standardRenderer *standard.Renderer
	shadowMap        *shadow.Renderer
	terrainRenderer  *terrain.Renderer
	skybox           *technique.SkyBox

	multiSampledTexture *framebuffer.Texture
	offScreenTexture    *framebuffer.Texture

	fullScreenTemp *framebuffer.Texture
}

func (e *Engine) SetActiveLight(light components.Light) {
	e.activeLight = light
}

func (e *Engine) ActiveLight() components.Light {
	return e.activeLight
}

func (e *Engine) Render(object, terrains components.Renderable) {
	if e.mainCamera == nil {
		panic("mainCamera not found, the game cannot render")
	}
	debugger.Clear()

	e.shadowMap.Render(object, terrains)

	e.multiSampledTexture.BindFrameBuffer()
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	e.shadowMap.Load()
	e.skybox.Load()
	e.terrainRenderer.Render(terrains)
	e.standardRenderer.Render(object)
	e.skybox.Render()
	e.multiSampledTexture.UnbindFrameBuffer()

	e.multiSampledTexture.ResolveToFBO(e.offScreenTexture)
	debug.CheckForError("renderer.Engine.Draw [end]")
	//e.offScreenTexture.ResolveToScreen()
	gl.Viewport(0, 0, e.width, e.height)
	gl.Disable(gl.DEPTH_TEST)
	e.applyFilter(e.toneMapShader, e.offScreenTexture, nil)
	//e.applyFilter(e.toneMapShader, e.offScreenTexture, e.fullScreenTemp)
	//e.applyFilter(e.overlayShader, debugger.Texture(), nil)
}

func (e *Engine) renderToScreen(t components.Texture) {
	e.applyFilter(e.nullShader, t, nil)
}

func (e *Engine) Lights() []components.Light {
	return e.lights
}

func (e *Engine) AddLight(l components.Light) {
	e.lights = append(e.lights, l)
}

func (e *Engine) applyFilter(shader components.Shader, in, out components.Texture) {
	if in == out {
		panic("Argh, can't apply shader where source and destination is the same")
	}

	if out == nil {
		gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
	} else {
		out.BindFrameBuffer()
	}
	e.SetTexture("x_filterTexture", in)
	e.SetInteger("x_w", in.Width())
	e.SetInteger("x_h", in.Height())
	shader.Bind()
	shader.UpdateUniforms(nil, e)
	primitives.DrawQuad()
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
	// @todo set value, regardless, this might be an array that isn't used
	v, ok := e.uniforms3f[name]
	if !ok {
		fmt.Printf("Vector3f, no value found for uniform '%s'\n", name)
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

func (e *Engine) SetSamplerSlot(samplerName string, slot uint32) {
	e.samplerMap[samplerName] = slot
}
