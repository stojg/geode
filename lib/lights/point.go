package lights

import (
	"github.com/go-gl/mathgl/mgl32"
)

func NewPoint(r, g, b, intensity float32) *PointLight {
	color := mgl32.Vec3{r, g, b}
	pointLight := &PointLight{
		BaseLight: BaseLight{
			color:     color.Mul(intensity),
			constant:  1,
			linear:    0.7,
			exponent:  1.80,
			intensity: intensity,
		},
	}

	pointLight.calcRange()

	return pointLight
}

type PointLight struct {
	BaseLight
}
