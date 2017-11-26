package lights

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/stojg/graphics/lib/components"
	"github.com/stojg/graphics/lib/rendering/framebuffer"
)

type BaseLight struct {
	components.GameComponent

	shadowInfo *ShadowInfo
	color      mgl32.Vec3
	shader     components.Shader
}

func (b *BaseLight) ShadowInfo() components.ShadowInfo {
	return b.shadowInfo
}

func (b *BaseLight) setShadowInfo(shadowInfo *ShadowInfo) {
	b.shadowInfo = shadowInfo
}

func (b *BaseLight) Shader() components.Shader {
	return b.shader
}

func (b *BaseLight) Color() mgl32.Vec3 {
	return b.color
}

func (b *BaseLight) SetColor(color mgl32.Vec3) {
	b.color = color
}

func (b *BaseLight) Position() mgl32.Vec3 {
	return b.Parent().Transform().Pos()
}

func (b *BaseLight) ShadowCaster() bool {
	return b.shadowInfo != nil
}

func (b *BaseLight) SetShadowTexture(slot uint32, texture *framebuffer.Texture) {}

func (b *BaseLight) BindShadow() {}

func (b *BaseLight) ViewProjection() mgl32.Mat4 {
	return mgl32.Ident4()
}

func (b *BaseLight) GetProjection() mgl32.Mat4 {
	return b.shadowInfo.projection
}

func (b *BaseLight) GetView() mgl32.Mat4 {
	//This comes from the conjugate rotation because the world should appear to rotate opposite to the camera's rotation.
	cameraRotation := b.Transform().TransformedRot().Conjugate().Mat4()
	//Similarly, the translation is inverted because the world appears to move opposite to the camera's movement.
	cameraPos := b.Transform().TransformedPos().Mul(-1)
	cameraTranslation := mgl32.Translate3D(cameraPos[0], cameraPos[1], cameraPos[2])
	return cameraRotation.Mul4(cameraTranslation)
}

func (b *BaseLight) Pos() mgl32.Vec3 {
	return b.Transform().Pos()
}
