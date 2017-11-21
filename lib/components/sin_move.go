package components

import (
	"time"

	"github.com/go-gl/mathgl/mgl32"
)

func NewTimeMove(axis mgl32.Vec3, f func(float64) float64) *SineMove {
	return &SineMove{
		axis: axis,
		f:    f,
	}
}

type SineMove struct {
	GameComponent
	axis         mgl32.Vec3
	totalElapsed time.Duration
	f            func(float64) float64
}

func (s *SineMove) Update(elapsed time.Duration) {
	s.totalElapsed += elapsed
	amount := float32(s.f(s.totalElapsed.Seconds()) * 0.003)
	s.Transform().MoveBy(s.axis.Mul(amount))
}
