package lights

import (
	"math"
	"time"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/stojg/graphics/lib/components"
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
	light.calcRange()
	return light
}

type Spot struct {
	BaseLight
	direction     mgl32.Vec3
	viewDirection mgl32.Mat4
}

func (spot *Spot) Update(elapsed time.Duration) {
	spot.BaseLight.Update(elapsed)
	spot.viewDirection = spot.Projection().Mul4(spot.calcView())
	spot.direction = spot.Transform().TransformedRot().Rotate(mgl32.Vec3{0, 0, -1})
}

func (spot *Spot) Direction() mgl32.Vec3 {
	return spot.direction
}

func (spot *Spot) AddToEngine(e components.RenderState) {
	e.AddLight(spot)
}

func (spot *Spot) ViewProjection() mgl32.Mat4 {
	return spot.viewDirection
}
