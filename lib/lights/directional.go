package lights

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/stojg/graphics/lib/components"
	"github.com/stojg/graphics/lib/rendering"
)

func NewDirectional(r, g, b, intensity float32) *Directional {
	const nearPlane float32 = 0.1
	const farPlane float32 = 25

	projection := mgl32.Ortho(-9, 9, -5, 18, nearPlane, farPlane)

	light := &Directional{
		BaseLight: BaseLight{
			color:      mgl32.Vec3{r, g, b}.Mul(intensity),
			shader:     rendering.NewShader("forward_directional"),
			shadowInfo: NewShadowInfo(projection, false),
		},
	}

	light.shadowInfo.shadowVarianceMin = 0.00002
	light.shadowInfo.lightBleedReduction = 0.8
	return light
}

type Directional struct {
	BaseLight

	view mgl32.Mat4
}

func (b *Directional) AddToEngine(e components.Engine) {
	e.GetRenderingEngine().AddLight(b)
}

func (b *Directional) Direction() mgl32.Vec3 {
	return b.BaseLight.Position().Normalize()
}

func (b *Directional) ViewProjection() mgl32.Mat4 {
	lightView := mgl32.LookAt(b.Position().X(), b.Position().Y(), b.Position().Z(), 0, 0, 0, 0, 1, 0)
	return b.shadowInfo.Projection().Mul4(lightView)
}
