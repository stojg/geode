package standard

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/stojg/geode/lib/components"
	"github.com/stojg/geode/lib/rendering/shader"
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
	object.Render(r.Camera(), r.shader, r, components.StandardRender)
}
