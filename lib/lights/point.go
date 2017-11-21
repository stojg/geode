package lights

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/stojg/graphics/lib/components"
	"github.com/stojg/graphics/lib/rendering"
)

func NewPoint(r, g, b, intensity float32) *PointLight {
	return &PointLight{
		BaseLight: BaseLight{
			color:  mgl32.Vec3{r, g, b}.Mul(intensity),
			shader: rendering.NewShader("forward_point"),
		},
		constant: 1,
		linear:   0.22,
		exponent: 0.20,
	}
}

type PointLight struct {
	BaseLight

	constant float32
	linear   float32
	exponent float32
}

func (p *PointLight) Exponent() float32 {
	return p.exponent
}

func (p *PointLight) Linear() float32 {
	return p.linear
}

func (p *PointLight) Constant() float32 {
	return p.constant
}

func (b *PointLight) AddToEngine(e components.Engine) {
	e.GetRenderingEngine().AddLight(b)
}
