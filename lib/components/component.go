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

type Light interface {
	Component
	SetShader(shader Shader)
	Shader() Shader
	Color() mgl32.Vec3
	SetColor(color mgl32.Vec3)
	Intensity() float32
	SetIntensity(intensity float32)
	Position() mgl32.Vec3
}

type RenderingEngine interface {
	AddLight(light Light)
	AddCamera(camera *Camera)
	GetMainCamera() *Camera
	GetActiveLight() Light
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

type GameComponent struct {
	parent Transformable
}

func (m *GameComponent) SetParent(parent Transformable) {
	m.parent = parent
}

func (m *GameComponent) Transform() *physics.Transform {
	return m.parent.Transform()
}

func (m *GameComponent) AddToEngine(engine Engine) {
}

func (m *GameComponent) Render(Shader, RenderingEngine) {}
func (m *GameComponent) Input(time.Duration)            {}
func (m *GameComponent) Update(time.Duration)           {}
