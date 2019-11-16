package geometry

import (
	"math"

	"github.com/go-gl/mathgl/mgl32"
)

type Plane [4]float32

func (p *Plane) Normalise() {
	l := float32(math.Sqrt(float64(p[0]*p[0] + p[1]*p[1] + p[2]*p[2])))
	p[0] /= l
	p[1] /= l
	p[2] /= l
	p[3] /= l
}

func (p *Plane) Dot(b mgl32.Vec3) float32 {
	x := p[0] * b[0]
	y := p[1] * b[1]
	z := p[2] * b[2]
	return x + y + z
}

func (p *Plane) DistanceToPoint(pt [3]float32, distance *float32) {
	a := p[0] * pt[0]
	b := p[1] * pt[1]
	c := p[2] * pt[2]
	*distance = a + b + c + p[3]
}
