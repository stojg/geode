package lights

import (
	"math"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/stojg/graphics/lib/components"
	"github.com/stojg/graphics/lib/rendering"
)

func NewSpot(r, g, b, intensity, angle float32) *Spot {

	c := float32(math.Cos(float64(mgl32.DegToRad(angle))))

	return &Spot{
		BaseLight: BaseLight{
			color:  mgl32.Vec3{r, g, b}.Mul(intensity),
			shader: rendering.NewShader("forward_spot"),
		},
		PointLight: PointLight{
			constant: 1,
			linear:   0.22,
			exponent: 0.20,
		},
		direction: mgl32.Vec3{0, 0, 0},
		cutoff:    c,
	}
}

type Spot struct {
	BaseLight
	PointLight

	direction mgl32.Vec3
	cutoff    float32
}

func (c *Spot) Direction() mgl32.Vec3 {
	r := c.Transform().TransformedRot()

	t := r.Rotate(mgl32.Vec3{0, 0, -1})
	//fmt.Println()
	return t
	//return mgl32.Vec3{0, -1, 0}
}

func (b *Spot) Cutoff() float32 {
	return b.cutoff
}

func (b *Spot) AddToEngine(e components.Engine) {
	e.GetRenderingEngine().AddLight(b)
}
