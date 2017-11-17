package components

import (
	"github.com/stojg/graphics/lib/physics"
)

type UniformUpdater interface {
}

type Drawable interface {
	Draw()
}

type Shader interface {
	Bind()
	UpdateUniforms(*physics.Transform, RenderingEngine)
}

type Transformable interface {
	Transform() *physics.Transform
}

type RenderingEngine interface {
	AddCamera(camera *Camera)
	GetMainCamera() *Camera
}

type Engine interface {
	GetRenderingEngine() RenderingEngine
}

type Component interface {
	Update(float32)
	Input(float32)
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
func (m *GameComponent) Input(float32)                  {}
func (m *GameComponent) Update(float32)                 {}
