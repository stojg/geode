package lights

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/stojg/graphics/lib/rendering"
)

func NewDirectional(r, g, b, intensity float32) *BaseLight {
	return &BaseLight{
		color:  mgl32.Vec3{r, g, b}.Mul(intensity),
		shader: rendering.NewShader("forward_directional"),
	}
}
