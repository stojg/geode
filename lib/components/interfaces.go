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

type ShadowCaster interface {
	BindAsRenderTarget()
	ShadowShader() Shader
	BindShadowTexture(samplerSlot uint32, samplerName string)
}

//type ShaderInfo interface {
//	Projection().
//}

type Light interface {
	Shader() Shader
	//ShadowInfo() ShaderInfo
	Color() mgl32.Vec3
	Position() mgl32.Vec3
	ViewProjection() mgl32.Mat4
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
