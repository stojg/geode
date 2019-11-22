package components

import (
	"time"

	"github.com/stojg/geode/lib/physics"
)

type GameComponent struct {
	parent Object
}

func (m *GameComponent) SetParent(parent Object) {
	m.parent = parent
}

func (m *GameComponent) Parent() Object {
	return m.parent
}

func (m *GameComponent) Transform() *physics.Transform {
	return m.parent.Transform()
}

func (m *GameComponent) AddToEngine(state RenderState) {
}

func (m *GameComponent) Update(time.Duration) {}
