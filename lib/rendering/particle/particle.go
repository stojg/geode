package particle

import (
	"time"

	"github.com/go-gl/mathgl/mgl32"
)

func NewParticle(pos, vel [3]float32, gravity, life, rotation, scale float32) *Particle {
	return &Particle{
		Position:   pos,
		Velocity:   vel,
		Gravity:    gravity,
		LifeLength: life,
		Rotation:   rotation,
		Scale:      scale,
	}
}

type Particle struct {
	Position   [3]float32
	Velocity   [3]float32
	Gravity    float32
	LifeLength float32
	Rotation   float32
	Scale      float32

	elapsedTime float32
}

func (p *Particle) Update(elapsed time.Duration) bool {
	t := float32(elapsed.Seconds())
	p.Velocity[1] += 50.0 * p.Gravity * t
	change := mgl32.Vec3(p.Position).Mul(t)
	p.Position[0] += change[0]
	p.Position[1] += change[1]
	p.Position[2] += change[2]
	p.elapsedTime += t
	return p.elapsedTime < p.LifeLength
}
