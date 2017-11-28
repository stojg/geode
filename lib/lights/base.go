package lights

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/stojg/graphics/lib/components"
)

type BaseLight struct {
	components.GameComponent
	shadowInfo *ShadowInfo
	color      mgl32.Vec3
	shader     components.Shader
}

func (b *BaseLight) AddToEngine(e components.Engine) {
	e.RenderingEngine().AddLight(b)
}

func (b *BaseLight) Shader() components.Shader {
	return b.shader
}

func (b *BaseLight) Color() mgl32.Vec3 {
	return b.color
}

func (b *BaseLight) Position() mgl32.Vec3 {
	return b.Parent().Transform().Pos()
}

func (b *BaseLight) ShadowInfo() components.ShadowInfo {
	return b.shadowInfo
}

func (b *BaseLight) ShadowCaster() bool {
	return b.shadowInfo != nil
}

func (b *BaseLight) ViewProjection() mgl32.Mat4 {
	return mgl32.Ident4()
}

func (b *BaseLight) Projection() mgl32.Mat4 {
	if b.shadowInfo == nil {
		return mgl32.Ident4()
	}
	return b.shadowInfo.projection
}

func (b *BaseLight) View() mgl32.Mat4 {
	//This comes from the conjugate rotation because the world should appear to rotate opposite to the camera's rotation.
	lightRotation := b.Transform().TransformedRot().Conjugate().Mat4()
	//Similarly, the translation is inverted because the world appears to move opposite to the camera's movement.
	lightPosition := b.Transform().TransformedPos().Mul(-1)
	lightTranslation := mgl32.Translate3D(lightPosition[0], lightPosition[1], lightPosition[2])
	return lightRotation.Mul4(lightTranslation)
}
