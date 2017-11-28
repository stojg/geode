package lights

import (
	"fmt"
	"math"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/stojg/graphics/lib/components"
	"github.com/stojg/graphics/lib/rendering"
)

func NewSpot(shadowSize int, r, g, b, intensity, viewAngle float32) *Spot {
	color := mgl32.Vec3{r, g, b}
	fov := mgl32.DegToRad(viewAngle)
	cutoff := float32(math.Cos(float64(fov / 2)))
	const nearPlane float32 = 2
	const farPlane float32 = 15

	light := &Spot{
		BaseLight: BaseLight{
			color:       color.Mul(intensity),
			shader:      rendering.NewShader("forward_spot"),
			maxDistance: 1,
		},
		PointLight: PointLight{
			constant: 1,
			linear:   0.35,
			exponent: 0.44,
		},
		cutoff: cutoff,
	}
	if shadowSize != 0 {
		projection := mgl32.Perspective(fov, float32(1), nearPlane, farPlane)
		light.shadowInfo = NewShadowInfo(shadowSize, projection, false)
		light.shadowInfo.shadowVarianceMin = 0.00002
		light.shadowInfo.lightBleedReduction = 0.8
	}

	max := color[0]
	if color[1] > max {
		max = color[1]
	}
	if color[2] > max {
		max = color[2]
	}

	fmt.Println(color)
	fmt.Println(max)

	const colorDepth = 256

	{
		a := light.Exponent()
		b := light.Linear()
		c := light.Constant() - colorDepth*intensity*max
		dist := (-b + float32(math.Sqrt(float64(b*b-4*a*c)))) / (2 * a)
		fmt.Println(dist)
		light.BaseLight.maxDistance = dist
	}

	return light
}

type Spot struct {
	BaseLight
	PointLight

	cutoff float32 // radians
}

func (spot *Spot) AddToEngine(e components.Engine) {
	e.RenderingEngine().AddLight(spot)
}

func (spot *Spot) Direction() mgl32.Vec3 {
	return spot.Transform().TransformedRot().Rotate(mgl32.Vec3{0, 0, -1})
}

func (spot *Spot) Cutoff() float32 {
	return spot.cutoff
}

func (spot *Spot) ViewProjection() mgl32.Mat4 {
	return spot.Projection().Mul4(spot.View())
}
