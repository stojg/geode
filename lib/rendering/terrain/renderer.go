package terrain

import (
	"github.com/stojg/graphics/lib/components"
	"github.com/stojg/graphics/lib/rendering/shader"
)

func NewRenderer(s components.RenderState) *Renderer {
	return &Renderer{
		RenderState: s,
		shader:      shader.NewShader("terrain"),
	}
}

type Renderer struct {
	components.RenderState
	shader     components.Shader
	mainCamera components.Viewable
	lights     []components.Light
	samplerMap map[string]uint32
}

func (r *Renderer) Render(object components.Renderable) {
	object.RenderAll(r.shader, r)
}
