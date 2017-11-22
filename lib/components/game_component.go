package components

import (
	"time"

	"github.com/stojg/graphics/lib/physics"
)

type GameComponent struct {
	parent Transformable
}

func (m *GameComponent) SetParent(parent Transformable) {
	m.parent = parent
}

func (m *GameComponent) Parent() Transformable {
	return m.parent
}

func (m *GameComponent) Transform() *physics.Transform {
	return m.parent.Transform()
}

func (m *GameComponent) AddToEngine(engine Engine) {
}

func (m *GameComponent) Render(Shader, RenderingEngine) {}
func (m *GameComponent) Input(time.Duration)            {}
func (m *GameComponent) Update(time.Duration)           {}
