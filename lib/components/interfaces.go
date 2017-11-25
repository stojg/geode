package components

import (
	"time"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/stojg/graphics/lib/physics"
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

type ShadowInfo interface {
	Projection() mgl32.Mat4
}

type Light interface {
	Shader() Shader
	Color() mgl32.Vec3
	Position() mgl32.Vec3
	ViewProjection() mgl32.Mat4
	ShadowInfo() ShadowInfo
	ShadowCaster() bool
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
	GetTexture(string) Texture
	SetTexture(string, Texture)

	GetVector3f(string) mgl32.Vec3
	SetVector3f(string, mgl32.Vec3)

	GetInteger(string) int32
	SetInteger(string, int32)
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
