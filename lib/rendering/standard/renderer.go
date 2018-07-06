package standard

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/stojg/graphics/lib/components"
	"github.com/stojg/graphics/lib/rendering/shader"
)

func NewRenderer(s components.RenderState) *Renderer {
	return &Renderer{
		RenderState: s,
		shader:      shader.NewShader("default"),
	}
}

type Renderer struct {
	components.RenderState
	shader components.Shader
}

func (r *Renderer) Render(object components.Renderable) {
	gl.Enable(gl.DEPTH_TEST)
	object.RenderAll(r.Camera(), r.shader, r)
}
