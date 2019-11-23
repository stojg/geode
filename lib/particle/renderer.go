package particle

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/stojg/geode/lib/components"
	"github.com/stojg/geode/lib/rendering/shader"
)

func NewRenderer(s components.RenderState) *Renderer {
	r := &Renderer{
		state:  s,
		shader: shader.NewShader("particle"),
	}

	return r
}

type Renderer struct {
	state  components.RenderState
	shader components.Shader
}

func (r *Renderer) Render(objects components.Renderable) {
	r.shader.Bind()
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.DepthMask(false)

	view := r.state.Camera().View()
	r.state.SetVector3f("x_camRight", [3]float32{view[0], view[4], view[8]})
	r.state.SetVector3f("x_camUp", [3]float32{view[1], view[5], view[9]})

	objects.Render(r.state.Camera(), r.shader, r.state, components.ParticleRender)

	gl.DepthMask(true)
	gl.Disable(gl.BLEND)
}
