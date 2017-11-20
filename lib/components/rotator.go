package components

import (
	"time"

	"github.com/go-gl/mathgl/mgl32"
)

func NewRotator(axis mgl32.Vec3, speed float32) *Rotator {
	return &Rotator{
		axis:  axis,
		speed: speed,
	}
}

type Rotator struct {
	GameComponent
	axis  mgl32.Vec3
	speed float32
}

func (r *Rotator) Update(elapsed time.Duration) {
	rot := float32(elapsed.Seconds()) * r.speed
	s := mgl32.DegToRad(rot)
	r.Transform().Rotate(r.axis, s)
}
