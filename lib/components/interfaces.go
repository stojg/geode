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
}

type Transformable interface {
	Transform() *physics.Transform
}

type Renderable interface {
	RenderAll(shader Shader, engine RenderingEngine)
}

type Viewable interface {
	View() mgl32.Mat4
	Projection() mgl32.Mat4
	Pos() mgl32.Vec3
	Rot() mgl32.Quat
}

type ShadowInfo interface {
	SizeAsPowerOfTwo() int
	Projection() mgl32.Mat4
	FlipFaces() bool
	LightBleedReduction() float32
	ShadowVarianceMin() float32
}

type Light interface {
	Shader() Shader
	Color() mgl32.Vec3
	MaxDistance() float32
	Position() mgl32.Vec3
	ViewProjection() mgl32.Mat4
	SetCamera(pos mgl32.Vec3, rot mgl32.Quat)
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
	AddCamera(camera Viewable)
	MainCamera() Viewable
	ActiveLight() Light
	SamplerSlot(string) uint32
	Texture(string) Texture
	SetTexture(string, Texture)

	Vector3f(string) mgl32.Vec3
	SetVector3f(string, mgl32.Vec3)

	Integer(string) int32
	SetInteger(string, int32)

	Float(string) float32
	SetFloat(string, float32)
}

type Engine interface {
	RenderingEngine() RenderingEngine
}

type Component interface {
	Update(time.Duration)
	Input(time.Duration)
	Render(Shader, RenderingEngine)
	AddToEngine(Engine)
	SetParent(Transformable)
}
