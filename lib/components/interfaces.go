package components

import (
	"time"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/stojg/geode/lib/geometry"
	"github.com/stojg/geode/lib/physics"
)

// The following constants are set on Gameobject and can be combines with binary | to decide which features they have
const (
	_          = iota // ignore first value by assigning to blank identifier
	ResourceNA = 1 << iota
	StandardRender
	ResourceTerrain
	ResourceLight
	ParticleRender
	Shadowed
)

type FBO interface {
	BindFrameBuffer()
	UnbindFrameBuffer()
}

type Texture interface {
	FBO

	ID() uint32
	Activate(samplerSlot uint32)

	Width() int32
	Height() int32
}

type Updatable interface {
	Update(time.Duration)
}

type Material interface {
	Texture(name string) Texture
	Textures() map[string]Texture
	AddTexture(name string, texture Texture)

	Float(name string) float32
	AddFloat(name string, float float32)

	Vector(name string) mgl32.Vec3
	AddVector(name string, vector mgl32.Vec3)
}

type Terrain interface {
	Height(x, z float32) float32
}

type Bindable interface {
	Bind()
}

type Unbindable interface {
	Unbind()
}

type AABB interface {
	C() mgl32.Vec3
	R() mgl32.Vec3
}

type Component interface {
	Updatable
	AddToEngine(state RenderState)
	SetParent(Object)
}

// ie mesh
type Drawable interface {
	AABB() AABB
	Bindable
	Unbindable
	Draw()
}

type Shader interface {
	Bindable
	Unbindable
	UpdateUniforms(Material Material, state RenderState)
	UpdateTransform(*physics.Transform, RenderState)
	UpdateUniform(name string, value interface{})
}

type Transformable interface {
	Transform() *physics.Transform
}

type Model interface {
	AABB() AABB
	Unbindable
	Updatable
	Bind(Shader, RenderState)
	Draw()
	Material() Material
}

type Object interface {
	Transformable
	Updatable
	Model() Model
	AllChildren() []Object
	IsType(int) bool
	Type() int
	IsVisible(camera Viewer) bool
	// @todo, the below method feels weird
	SetState(state RenderState)
	Draw(camera Viewer, shader Shader, state RenderState)
}

type Renderable interface {
	Render(camera Viewer, shader Shader, state RenderState, rtype int)
}

type Viewer interface {
	Frustum() geometry.Frustum
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

type RenderState interface {
	Update()

	Camera() Viewer
	SetCamera(camera Viewer)

	SamplerSlot(name string) uint32
	AddSamplerSlot(name string)

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

type Renderer interface {
	Render(a Renderable)
	State() RenderState
}

type Engine interface {
	Renderer() Renderer
}

type Logger interface {
	Println(a ...interface{})
	Printf(format string, a ...interface{})
	ErrorLn(inError error)
}
