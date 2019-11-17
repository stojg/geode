package particle

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/stojg/geode/lib/components"
	"github.com/stojg/geode/lib/rendering/shader"
)

func NewRenderer(s components.RenderState) *Renderer {
	r := &Renderer{
		RenderState: s,
		shader:      shader.NewShader("particle"),
	}

	return r
}

type Renderer struct {
	components.RenderState
	shader components.Shader
}

func (r *Renderer) Render(objects components.Renderable) {
	r.shader.Bind()
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.DepthMask(false)

	objects.Render(r.RenderState.Camera(), r.shader, r.RenderState, components.ParticleRender)

	gl.DepthMask(true)
	gl.Disable(gl.BLEND)
}
