package lights

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/stojg/graphics/lib/components"
	"github.com/stojg/graphics/lib/rendering/framebuffer"
)

func NewShadowInfo(projection mgl32.Mat4) *ShadowInfo {
	return &ShadowInfo{
		projection: projection,
	}
}

type ShadowInfo struct {
	projection mgl32.Mat4
}

func (s *ShadowInfo) Projection() mgl32.Mat4 {
	return s.projection
}

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
