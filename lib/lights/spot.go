package lights

import (
	"math"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/stojg/graphics/lib/components"
	"github.com/stojg/graphics/lib/rendering"
)

func NewSpot(r, g, b, intensity, viewAngle float32) *Spot {

	fov := mgl32.DegToRad(viewAngle)
	cutoff := float32(math.Cos(float64(fov / 2)))
	const nearPlane float32 = 2
	const farPlane float32 = 15

	projection := mgl32.Perspective(fov, float32(1), nearPlane, farPlane)

	light := &Spot{
		BaseLight: BaseLight{
			color:      mgl32.Vec3{r, g, b}.Mul(intensity),
			shader:     rendering.NewShader("forward_spot"),
			shadowInfo: NewShadowInfo(projection, false),
		},
		PointLight: PointLight{
			constant: 1,
			linear:   0.22,
			exponent: 0.20,
		},
		cutoff: cutoff,
	}
	light.shadowInfo.shadowVarianceMin = 0.00002
	light.shadowInfo.lightBleedReduction = 0.8
	return light
}

type Spot struct {
	BaseLight
	PointLight
	cutoff float32 // radians
}

func (spot *Spot) AddToEngine(e components.Engine) {
	e.GetRenderingEngine().AddLight(spot)
}

func (spot *Spot) Direction() mgl32.Vec3 {
	return spot.Transform().TransformedRot().Rotate(mgl32.Vec3{0, 0, -1})
}

func (spot *Spot) Cutoff() float32 {
	return spot.cutoff
}

func (spot *Spot) ViewProjection() mgl32.Mat4 {
	return spot.GetProjection().Mul4(spot.GetView())
}
