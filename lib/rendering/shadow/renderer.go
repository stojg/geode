package shadow

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/stojg/geode/lib/components"
	"github.com/stojg/geode/lib/rendering/framebuffer"
	"github.com/stojg/geode/lib/rendering/shader"
)

func NewRenderer(s components.RenderState) *Renderer {
	e := &Renderer{
		RenderState: s,
		gaussShader: shader.NewShader("filter_gauss"),
		shader:      shader.NewShader("shadow_vsm"),
	}

	e.shadowTextures = make([]components.Texture, 12)
	e.tempShadowTextures = make([]components.Texture, 12)
	for i := uint(0); i < 12; i++ {
		size := 1 << i // power of two, 1, 2, 4, 8, 16 and so on
		e.shadowTextures[i] = framebuffer.NewTexture(gl.COLOR_ATTACHMENT0, size, size, gl.RG32F, gl.RG, gl.FLOAT, gl.LINEAR, true)
		e.tempShadowTextures[i] = framebuffer.NewTexture(gl.COLOR_ATTACHMENT0, size, size, gl.RG32F, gl.RG, gl.FLOAT, gl.LINEAR, true)
	}
	s.AddSamplerSlot("x_shadowMap")

	return e
}

type Renderer struct {
	components.RenderState
	gaussShader        components.Shader
	shader             components.Shader
	shadowTextures     []components.Texture
	tempShadowTextures []components.Texture
	shadowCaster       components.Light
}

func (r *Renderer) Render(objects components.Renderable) {
	gl.Enable(gl.DEPTH_TEST)
	for _, light := range r.RenderState.Lights() {
		if light.ShadowCaster() {
			r.shadowCaster = light
			break
		}
	}
	r.SetActiveLight(r.shadowCaster)
	r.shadowCaster.SetCamera(r.Camera().Pos(), r.Camera().Rot())
	idx := r.shadowCaster.ShadowInfo().SizeAsPowerOfTwo()
	r.shadowTextures[idx].BindFrameBuffer()
	gl.ClearColor(1, 1, 1, 1)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	objects.Render(r.Camera(), r.shader, r, components.R_SHADOWED)
	//terrains.Render(r.Camera(), r.shader, r)
	//e.blurShadowMap(idx, 1)
}

func (r *Renderer) Load() {
	r.RenderState.SetFloat("x_varianceMin", r.ShadowVarianceMin())
	r.RenderState.SetFloat("x_lightBleedReductionAmount", r.LightBleedReduction())
	r.RenderState.SetTexture("x_shadowMap", r.ShadowMapTexture())
}

func (r *Renderer) ShadowVarianceMin() float32 {
	return r.shadowCaster.ShadowInfo().ShadowVarianceMin()
}
func (r *Renderer) LightBleedReduction() float32 {
	return r.shadowCaster.ShadowInfo().LightBleedReduction()
}

func (r *Renderer) ShadowMapTexture() components.Texture {
	return r.shadowTextures[r.shadowCaster.ShadowInfo().SizeAsPowerOfTwo()]
}

//func (e *Renderer) blurShadowMap(sizeAsPowerOfTwo int, blurAmount float32) {
//	var size = 2 << uint(sizeAsPowerOfTwo)
//	src := e.shadowTextures[sizeAsPowerOfTwo]
//	tmp := e.tempShadowTextures[sizeAsPowerOfTwo]
//	gl.Disable(gl.DEPTH_TEST)
//	defer gl.Enable(gl.DEPTH_TEST)
//	gl.Viewport(0, 0, src.Width(), src.Height())
//	e.SetVector3f("x_blurScale", mgl32.Vec3{1 / float32(size) * blurAmount, 0, 0})
//	e.applyFilter(e.gaussShader, src, tmp)
//	e.SetVector3f("x_blurScale", mgl32.Vec3{0, 1 / float32(size) * blurAmount, 0})
//	e.applyFilter(e.gaussShader, tmp, src)
//	gl.GenerateMipmap(gl.TEXTURE_2D)
//
//}
