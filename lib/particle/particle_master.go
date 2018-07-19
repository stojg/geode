package particle

import (
	"github.com/stojg/graphics/lib/components"
)

func NewMaster(s components.RenderState) *Master {
	return &Master{
		renderer: NewRenderer(s),
	}

}

type Master struct {
	particles []*Particle
	renderer  *Renderer
}

func (m *Master) Update(elapsed float32) {
	for i := len(m.particles) - 1; i >= 0; i-- {
		alive := m.particles[i].Update(elapsed)
		if !alive {
			m.particles = append(m.particles[:i], m.particles[i+1:]...)
		}
	}
}

func (m *Master) Render(camera components.Viewable) {
	m.renderer.Render(m.particles, camera)
}

func (m *Master) AddParticle(pos, vel, rotAxis [3]float32, scale, rotAngle, gravity, life float32) {
	p := NewParticle(pos, vel, rotAxis, scale, rotAngle, gravity, life)
	m.particles = append(m.particles, p)
}

func (m *Master) Len() int {
	return len(m.particles)
}
