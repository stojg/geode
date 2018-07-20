package particle

import (
	"math/rand"
	"time"

	"github.com/stojg/graphics/lib/components"
)

func NewSystem() *System {

	return &System{}
}

type System struct {
	components.GameComponent
	particles []*Particle
}

func (s *System) Update(elapsed time.Duration) {
	secs := float32(elapsed.Seconds())
	for i := len(s.particles) - 1; i >= 0; i-- {
		alive := s.particles[i].Update(secs)
		if !alive {
			s.particles = append(s.particles[:i], s.particles[i+1:]...)
		}
	}

	if len(s.particles) < 100 {
		s.AddParticle([3]float32{0, 3.8, 0}, [3]float32{rand.Float32()*4 - 2, rand.Float32() * 20, rand.Float32()*4 - 2}, rand.Float32()*0.05+0.025, rand.Float32()*45, 1, rand.Float32()*10)
	}
}

func (m *System) AddParticle(pos, vel [3]float32, scale, rotAngle, gravity, life float32) {
	p := NewParticle(pos, vel, scale, rotAngle, gravity, life)
	m.particles = append(m.particles, p)
}
