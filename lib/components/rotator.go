package components

import (
	"time"

	"github.com/go-gl/mathgl/mgl32"
)

type Rotator struct {
	GameComponent
}

func (r *Rotator) Update(elapsed time.Duration) {
	rot := elapsed.Seconds() * 45 // 45 degrees per second
	s := mgl32.DegToRad(float32(rot))
	r.Transform().Rotate(vec3(1, 1, 1), s)
}
