package lights

import (
	"math"

	"github.com/go-gl/mathgl/mgl32"
)

func NewSpot(r, g, b, intensity, viewAngle float32) *Spot {
	color := mgl32.Vec3{r, g, b}
	fov := mgl32.DegToRad(viewAngle)
	cutoff := float32(math.Cos(float64(fov / 2)))

	light := &Spot{
		BaseLight: BaseLight{
			color:       color.Mul(intensity),
			maxDistance: 1,
			intensity:   intensity,
			constant:    1,
			linear:      0.35,
			exponent:    0.44,
			cutoff:      cutoff,
		},
	}
	calcRange(&light.BaseLight)
	return light
}

type Spot struct {
	BaseLight
}

func (spot *Spot) Direction() mgl32.Vec3 {
	return spot.Transform().TransformedRot().Rotate(mgl32.Vec3{0, 0, -1})
}

func (spot *Spot) ViewProjection() mgl32.Mat4 {
	return spot.Projection().Mul4(spot.View())
}
