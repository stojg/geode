package lights

import (
	"math"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/stojg/graphics/lib/components"
)

type BaseLight struct {
	components.GameComponent
	shadowInfo *ShadowInfo
	color      mgl32.Vec3

	maxDistance float32

	constant float32
	linear   float32
	exponent float32

	cutoff float32

	intensity float32
}

func (l *BaseLight) MaxDistance() float32 {
	return l.maxDistance
}

func (l *BaseLight) SetCamera(pos mgl32.Vec3, rot mgl32.Quat) {}

func (l *BaseLight) AddToEngine(e components.Engine) {
	e.Renderer().State().AddLight(l)
}

func (l *BaseLight) Direction() mgl32.Vec3 {
	return l.Parent().Transform().TransformedPos().Normalize()
}

func (l *BaseLight) Color() mgl32.Vec3 {
	return l.color
}

func (l *BaseLight) Exponent() float32 {
	return l.exponent
}

func (l *BaseLight) Linear() float32 {
	return l.linear
}

func (l *BaseLight) Constant() float32 {
	return l.constant
}

func (l *BaseLight) Cutoff() float32 {
	return l.cutoff
}

func (l *BaseLight) Position() mgl32.Vec3 {
	return l.Parent().Transform().Pos()
}

func (l *BaseLight) ShadowInfo() components.ShadowInfo {
	return l.shadowInfo
}

func (l *BaseLight) ShadowCaster() bool {
	return l.shadowInfo != nil
}

func (l *BaseLight) ViewProjection() mgl32.Mat4 {
	return mgl32.Ident4()
}

func (l *BaseLight) Projection() mgl32.Mat4 {
	if l.shadowInfo == nil {
		return mgl32.Ident4()
	}
	return l.shadowInfo.projection
}

func calcRange(l *BaseLight) {

	max := l.color[0]
	if l.color[1] > max {
		max = l.color[1]
	}
	if l.color[2] > max {
		max = l.color[2]
	}

	const colorDepth = 256

	{
		a := l.Exponent()
		b := l.Linear()
		c := l.Constant() - colorDepth*l.intensity*max
		dist := (-b + float32(math.Sqrt(float64(b*b-4*a*c)))) / (2 * a)
		l.maxDistance = dist
	}
}

func (l *BaseLight) View() mgl32.Mat4 {
	//This comes from the conjugate rotation because the world should appear to rotate opposite to the camera's rotation.
	lightRotation := l.Transform().TransformedRot().Conjugate().Mat4()
	//Similarly, the translation is inverted because the world appears to move opposite to the camera's movement.
	lightPosition := l.Transform().TransformedPos().Mul(-1)
	lightTranslation := mgl32.Translate3D(lightPosition[0], lightPosition[1], lightPosition[2])
	return lightRotation.Mul4(lightTranslation)
}
