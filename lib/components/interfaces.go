package components

import (
	"time"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/stojg/graphics/lib/physics"
)

type Texture interface {
	ID() uint32
	Activate(samplerSlot uint32)
	BindFrameBuffer()
	UnbindFrameBuffer()

	Width() int32
	Height() int32
}

type Material interface {
	Texture(name string) Texture
	Albedo() mgl32.Vec3
	Metallic() float32
	Roughness() float32
}

type Terrain interface {
	Height(x, z float32) float32
}

// ie mesh
type Drawable interface {
	Bind()
	Draw()
	Unbind()
	HalfWidths() [3][2]float32
}

type Logger interface {
	Println(a ...interface{})
	Printf(format string, a ...interface{})
	ErrorLn(inError error)
}

type Shader interface {
	Bind()
	UpdateUniforms(Material, RenderingEngine)
	UpdateTransform(*physics.Transform, RenderingEngine)
	UpdateUniform(name string, value interface{})
}

type Transformable interface {
	Transform() *physics.Transform
	AllTransforms() []*physics.Transform
}

type Renderable interface {
	RenderAll(camera Viewable, shader Shader, engine RenderingEngine)
}

type Viewable interface {
	Planes() [6][4]float32
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
	Color() mgl32.Vec3
	MaxDistance() float32
	Exponent() float32
	Linear() float32
	Constant() float32
	Cutoff() float32
	Direction() mgl32.Vec3

	Position() mgl32.Vec3
	ViewProjection() mgl32.Mat4
	SetCamera(pos mgl32.Vec3, rot mgl32.Quat)
	ShadowInfo() ShadowInfo
	ShadowCaster() bool
}

type DirectionalLight interface {
	Light
}

type PointLight interface {
	Light
}

type Spotlight interface {
	PointLight
}

type RenderState interface {
	AddCamera(camera Viewable)
	MainCamera() Viewable

	SamplerSlot(name string) uint32
	SetSamplerSlot(name string, slot uint32)

	AddLight(light Light)
	Lights() []Light
	SetActiveLight(light Light)
	ActiveLight() Light

	Texture(string) Texture
	SetTexture(string, Texture)

	Vector3f(string) mgl32.Vec3
	SetVector3f(string, mgl32.Vec3)

	Integer(string) int32
	SetInteger(string, int32)

	Float(string) float32
	SetFloat(string, float32)
}

type RenderingEngine interface {
	RenderState
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
