package components

import (
	"time"

	"github.com/go-gl/mathgl/mgl32"
)

func NewSpinner(axis mgl32.Vec3, velocity float32) *Spinner {
	return &Spinner{
		axis:     axis,
		velocity: velocity,
	}
}

type Spinner struct {
	GameComponent
	axis     mgl32.Vec3
	velocity float32
}

func (r *Spinner) Update(elapsed time.Duration) {
	rot := float32(elapsed.Seconds()) * r.velocity
	s := mgl32.DegToRad(rot)
	r.Transform().Rotate(r.axis, s)
}
