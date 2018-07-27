package postprocess

import (
	"fmt"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/stojg/graphics/lib/components"
	"github.com/stojg/graphics/lib/rendering/framebuffer"
	"github.com/stojg/graphics/lib/rendering/primitives"
	"github.com/stojg/graphics/lib/rendering/shader"
)

func New(s components.RenderState, width, height, vpWidth, vpHeight int) *Renderer {
	s.AddSamplerSlot("x_filterTexture")
	r := &Renderer{
		RenderState:       s,
		sourceTexture:     framebuffer.NewTexture(gl.COLOR_ATTACHMENT0, width, height, gl.RGBA16F, gl.RGBA, gl.FLOAT, gl.LINEAR, false),
		brightPassTexture: framebuffer.NewTexture(gl.COLOR_ATTACHMENT0, width, height, gl.RGBA16F, gl.RGBA, gl.FLOAT, gl.LINEAR, false),
		scratch1:          framebuffer.NewTexture(gl.COLOR_ATTACHMENT0, width, height, gl.RGBA16F, gl.RGBA, gl.FLOAT, gl.LINEAR, false),
		scratch2:          framebuffer.NewTexture(gl.COLOR_ATTACHMENT0, width, height, gl.RGBA16F, gl.RGBA, gl.FLOAT, gl.LINEAR, false),

		vpWidth:       vpWidth,
		vpHeight:      vpHeight,
		toneMapShader: shader.NewShader("filter_tonemap"),
		gaussShader:   shader.NewShader("filter_gauss"),
		brightness:    shader.NewShader("filter_brightness"),
		combine:       shader.NewShader("filter_combine"),
		pass:          shader.NewShader("filter_null"),
	}

	//r.pixelShader = shader.NewShader("filter_pixelator")

	for i := uint(2); i < 5; i++ {
		size := 1 << i // power of two, 1, 2, 4, 8, 16 and so on
		s.AddSamplerSlot(fmt.Sprintf("x_filterTexture%d", i))
		blurTextures := [2]*framebuffer.Texture{
			framebuffer.NewTexture(gl.COLOR_ATTACHMENT0, width/size, height/size, gl.RGB, gl.RGB, gl.HALF_FLOAT, gl.LINEAR, false),
			framebuffer.NewTexture(gl.COLOR_ATTACHMENT0, width/size, height/size, gl.RGB, gl.RGB, gl.HALF_FLOAT, gl.LINEAR, false),
		}
		r.blurTextures = append(r.blurTextures, blurTextures)
	}
	return r
}

type Renderer struct {
	vpWidth, vpHeight int
	components.RenderState
	sourceTexture     *framebuffer.Texture
	brightPassTexture *framebuffer.Texture
	scratch1          *framebuffer.Texture
	scratch2          *framebuffer.Texture
	blurTextures      [][2]*framebuffer.Texture
	brightnessTex     *framebuffer.Texture
	toneMapShader     *shader.Shader
	gaussShader       *shader.Shader
	brightness        *shader.Shader
	combine           *shader.Shader
	pass              *shader.Shader
	pixelShader       *shader.Shader
}

func (r *Renderer) Render(in *framebuffer.Texture, bypass bool) {

	gl.Disable(gl.DEPTH_TEST)
	if bypass {
		in.ResolveToFBO(r.sourceTexture)
		r.applyFilter(r.toneMapShader, r.sourceTexture, nil)
		return
	}

	in.ResolveToFBO(r.sourceTexture)

	r.applyFilter(r.brightness, r.sourceTexture, r.brightPassTexture)

	r.gaussShader.Bind()
	res := r.brightPassTexture
	for i, t := range r.blurTextures {
		res = r.blur(res, t[0], t[1])
		r.SetTexture(fmt.Sprintf("x_filterTexture%d", i+2), res)
	}
	r.applyFilter(r.combine, r.sourceTexture, r.scratch1)

	if r.pixelShader != nil {
		r.applyFilter(r.pixelShader, r.scratch1, r.scratch2)
		r.applyFilter(r.toneMapShader, r.scratch2, nil)
	} else {
		r.applyFilter(r.toneMapShader, r.scratch1, nil)
	}
}

func (r *Renderer) blur(in, t1, t2 *framebuffer.Texture) *framebuffer.Texture {
	initial := true
	x := true
	from, to := t1, t2
	for i := 0; i < 4; i++ {
		var scale [3]float32
		if x {
			scale = [3]float32{1, 0, 0}
		} else {
			scale = [3]float32{0, 1, 0}
		}
		r.SetVector3f("x_blurScale", scale)
		r.SetTexture("x_filterTexture", from)
		if initial {
			r.SetTexture("x_filterTexture", in)
			initial = false
		}
		to.BindFrameBuffer()
		r.gaussShader.UpdateUniforms(nil, r)
		primitives.DrawQuad()
		x = !x
		from, to = to, from
	}
	return from
}

func (r *Renderer) applyFilter(shader components.Shader, in, out components.Texture) {
	if out == nil {
		gl.Viewport(0, 0, int32(r.vpWidth), int32(r.vpHeight))
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
