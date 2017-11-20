package lights

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/stojg/graphics/lib/rendering"
)

func NewPoint(r, g, b, intensity float32) *PointLight {
	return &PointLight{
		BaseLight: BaseLight{
			color:  mgl32.Vec3{r, g, b}.Mul(intensity),
			shader: rendering.NewShader("forward_point"),
		},
		constant:  1,
		linear:    0.22,
		quadratic: 0.20,
	}
}

type PointLight struct {
	BaseLight

	constant  float32
	linear    float32
	quadratic float32
}
