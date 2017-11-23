package components

import (
	"time"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/stojg/graphics/lib/physics"
	"github.com/stojg/graphics/lib/rendering/framebuffer"
)

type Drawable interface {
	Draw()
}

type Shader interface {
	Bind()
	UpdateUniforms(*physics.Transform, Material, RenderingEngine)
	SetUniform(string, interface{})
}

type Transformable interface {
	Transform() *physics.Transform
}

type Renderable interface {
	RenderAll(shader Shader, engine RenderingEngine)
}

type ShadowCaster interface {
	SetShadowTexture(slot uint32, samplerName string, texture *framebuffer.Texture)
	BindShadow()
}

type Light interface {
	Shader() Shader
	Color() mgl32.Vec3
	Position() mgl32.Vec3
}

type DirectionalLight interface {
	Light
	Direction() mgl32.Vec3
}

type PointLight interface {
	Light
	Exponent() float32
	Linear() float32
	Constant() float32
}

type Spotlight interface {
	PointLight
	Direction() mgl32.Vec3
	Cutoff() float32
}

type RenderingEngine interface {
	AddLight(light Light)
	AddCamera(camera *Camera)
	GetMainCamera() *Camera
	GetActiveLight() Light
	GetSamplerSlot(string) uint32
}

type Engine interface {
	GetRenderingEngine() RenderingEngine
}

type Component interface {
	Update(time.Duration)
	Input(time.Duration)
	Render(Shader, RenderingEngine)
	AddToEngine(Engine)
	SetParent(Transformable)
}
