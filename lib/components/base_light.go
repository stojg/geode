package components

import "github.com/go-gl/mathgl/mgl32"

func NewBaseLight(color mgl32.Vec3, intensity float32) *BaseLight {
	return &BaseLight{
		color:     color,
		intensity: intensity,
	}
}

type BaseLight struct {
	GameComponent

	color     mgl32.Vec3
	intensity float32
	shader    Shader
}

func (b *BaseLight) AddToEngine(e Engine) {
	e.GetRenderingEngine().AddLight(b)
}

func (b *BaseLight) SetShader(shader Shader) {
	b.shader = shader
}

func (b *BaseLight) Shader() Shader {
	return b.shader
}

func (b *BaseLight) Color() mgl32.Vec3 {
	return b.color
}

func (b *BaseLight) SetColor(color mgl32.Vec3) {
	b.color = color
}

func (b *BaseLight) Intensity() float32 {
	return b.intensity
}

func (b *BaseLight) SetIntensity(intensity float32) {
	b.intensity = intensity
}

func (b *BaseLight) Position() mgl32.Vec3 {
	return b.parent.Transform().Pos()
}
