package postprocess

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/stojg/graphics/lib/components"
	"github.com/stojg/graphics/lib/rendering/framebuffer"
	"github.com/stojg/graphics/lib/rendering/primitives"
	"github.com/stojg/graphics/lib/rendering/shader"
)

func New(s components.RenderState) *Renderer {

	const blurDownScale = 4
	return &Renderer{
		RenderState:   s,
		sourceTexture: framebuffer.NewTexture(gl.COLOR_ATTACHMENT0, components.Width, components.Height, gl.RGBA16F, gl.RGBA, gl.FLOAT, gl.NEAREST, false),
		scratch1:      framebuffer.NewTexture(gl.COLOR_ATTACHMENT0, components.Width/2, components.Height/2, gl.RGBA16F, gl.RGB, gl.FLOAT, gl.NEAREST, false),
		scratch2:      framebuffer.NewTexture(gl.COLOR_ATTACHMENT0, components.Width, components.Height, gl.RGBA16F, gl.RGB, gl.FLOAT, gl.NEAREST, false),
		blur1:         framebuffer.NewTexture(gl.COLOR_ATTACHMENT0, components.Width/blurDownScale, components.Height/blurDownScale, gl.RGB, gl.RGB, gl.FLOAT, gl.LINEAR, false),
		blur2:         framebuffer.NewTexture(gl.COLOR_ATTACHMENT0, components.Width/blurDownScale, components.Height/blurDownScale, gl.RGB, gl.RGB, gl.FLOAT, gl.LINEAR, false),
		blur3:         framebuffer.NewTexture(gl.COLOR_ATTACHMENT0, components.Width/blurDownScale/2, components.Height/blurDownScale/2, gl.RGB, gl.RGB, gl.FLOAT, gl.LINEAR, false),
		blur4:         framebuffer.NewTexture(gl.COLOR_ATTACHMENT0, components.Width/blurDownScale/2, components.Height/blurDownScale/2, gl.RGB, gl.RGB, gl.FLOAT, gl.LINEAR, false),
		blur5:         framebuffer.NewTexture(gl.COLOR_ATTACHMENT0, components.Width/blurDownScale/4, components.Height/blurDownScale/4, gl.RGB, gl.RGB, gl.FLOAT, gl.LINEAR, false),
		blur6:         framebuffer.NewTexture(gl.COLOR_ATTACHMENT0, components.Width/blurDownScale/4, components.Height/blurDownScale/4, gl.RGB, gl.RGB, gl.FLOAT, gl.LINEAR, false),

		toneMapShader: shader.NewShader("filter_tonemap"),
		gaussShader:   shader.NewShader("filter_gauss"),
		brightness:    shader.NewShader("filter_brightness"),
		combine:       shader.NewShader("filter_combine"),
		pass:          shader.NewShader("filter_null"),
	}
}

type Renderer struct {
	components.RenderState
	sourceTexture *framebuffer.Texture
	scratch1      *framebuffer.Texture
	scratch2      *framebuffer.Texture
	blur1         *framebuffer.Texture
	blur2         *framebuffer.Texture
	blur3         *framebuffer.Texture
	blur4         *framebuffer.Texture
	blur5         *framebuffer.Texture
	blur6         *framebuffer.Texture
	brightnessTex *framebuffer.Texture
	toneMapShader *shader.Shader
	gaussShader   *shader.Shader
	brightness    *shader.Shader
	combine       *shader.Shader
	pass          *shader.Shader
}

func (r *Renderer) Render(in *framebuffer.Texture) {

	in.ResolveToFBO(r.sourceTexture)
	gl.Disable(gl.DEPTH_TEST)
	//gl.Viewport(0, 0, e.width, e.height)

	r.applyFilter(r.brightness, r.sourceTexture, r.scratch1)
	res1 := r.blur(r.scratch1, r.blur1, r.blur2)
	res2 := r.blur(res1, r.blur3, r.blur4)
	res3 := r.blur(res2, r.blur5, r.blur6)

	r.SetTexture("x_filterTexture2", res1)
	r.SetTexture("x_filterTexture3", res2)
	r.SetTexture("x_filterTexture4", res3)
	r.applyFilter(r.combine, r.sourceTexture, r.scratch2)
	r.applyFilter(r.toneMapShader, r.scratch2, nil)
}

func (r *Renderer) blur(in, t1, t2 *framebuffer.Texture) *framebuffer.Texture {
	initial := true
	x := true
	from, to := t1, t2
	for i := 0; i < 2; i++ {
		var scale [3]float32
		if x {
			scale = [3]float32{1, 0, 0}
		} else {
			scale = [3]float32{0, 1, 0}
		}
		r.SetVector3f("x_blurScale", scale)
		if initial {
			r.applyFilter(r.gaussShader, in, to)
			initial = false
		} else {
			r.applyFilter(r.gaussShader, from, to)
		}
		x = !x
		from, to = to, from
	}
	return from
}

func (r *Renderer) applyFilter(shader components.Shader, in, out components.Texture) {
	if out == nil {
		gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
	} else {
		out.BindFrameBuffer()
	}
	r.SetTexture("x_filterTexture", in)
	r.SetInteger("x_w", in.Width())
	r.SetInteger("x_h", in.Height())
	shader.Bind()
	shader.UpdateUniforms(nil, r)
	primitives.DrawQuad()
}
