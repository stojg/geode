package lights

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/stojg/graphics/lib/components"
	"github.com/stojg/graphics/lib/rendering"
)

func NewDirectional(r, g, b, intensity float32) *Directional {
	return &Directional{
		BaseLight: BaseLight{
			color:  mgl32.Vec3{r, g, b}.Mul(intensity),
			shader: rendering.NewShader("forward_directional"),
		},
	}
}

type Directional struct {
	BaseLight
}

func (b *Directional) AddToEngine(e components.Engine) {
	e.GetRenderingEngine().AddLight(b)
}
