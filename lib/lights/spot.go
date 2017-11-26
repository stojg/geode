package lights

import (
	"math"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/stojg/graphics/lib/components"
	"github.com/stojg/graphics/lib/rendering"
)

func NewSpot(r, g, b, intensity, angle float32) *Spot {

	fov := mgl32.DegToRad(angle)
	radians := float32(math.Cos(float64(fov)))
	const nearPlane float32 = 0.01
	const farPlane float32 = 20

	projection := mgl32.Ortho(-9, 9, -5, 18, nearPlane, farPlane)

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
		direction: mgl32.Vec3{0, 0, 0},
		cutoff:    radians,
	}
	light.shadowInfo.shadowVarianceMin = 0.00002
	light.shadowInfo.lightBleedReduction = 0.8
	return light
}

type Spot struct {
	BaseLight
	PointLight

	direction mgl32.Vec3
	// radians
	cutoff float32
}

func (c *Spot) Direction() mgl32.Vec3 {
	r := c.Transform().TransformedRot()

	t := r.Rotate(mgl32.Vec3{0, 0, -1})
	return t
}

func (b *Spot) Cutoff() float32 {
	return b.cutoff
}

func (b *Spot) AddToEngine(e components.Engine) {
	e.GetRenderingEngine().AddLight(b)
}
func (b *Spot) ViewProjection() mgl32.Mat4 {
	lightView := mgl32.LookAt(b.Position().X(), b.Position().Y(), b.Position().Z(), 0, 0, 0, 0, 1, 0)
	return b.shadowInfo.Projection().Mul4(lightView)
}
