package particle

import (
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/stojg/graphics/lib/components"
	"github.com/stojg/graphics/lib/rendering/shader"
)

var Vertices = [8]float32{-0.5, 0.5, -0.5, -0.5, 0.5, 0.5, 0.5, -0.5}

func NewRenderer(s components.RenderState) *Renderer {

	//quad := primitives.DrawQuad()

	r := &Renderer{
		RenderState: s,
		shader:      shader.NewShader("particle"),
	}
	return r
}

type Renderer struct {
	components.RenderState
	shader  components.Shader
	quadVao uint32
}

func (r *Renderer) Render(object components.Renderable) {
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.DepthMask(false)

	object.RenderAll(r.Camera(), r.shader, r)
}
