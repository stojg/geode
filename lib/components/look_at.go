package components

import (
	"time"

	"github.com/go-gl/mathgl/mgl32"
)

func NewLookAt(target mgl32.Vec3) *LookAt {
	return &LookAt{
		target: target,
	}
}

type LookAt struct {
	GameComponent
	target mgl32.Vec3
}

func (r *LookAt) Update(elapsed time.Duration) {
	r.Transform().LookAt(r.target, mgl32.Vec3{0, 1, 0})
}
