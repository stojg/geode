package particle

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/stojg/graphics/lib/components"
	"github.com/stojg/graphics/lib/math"
)

const Gravity float32 = -9.82

func NewParticle(pos, vel [3]float32, scale, rotAngle, gravity, life float32) *Particle {
	return &Particle{
		Position:     pos,
		Velocity:     vel,
		Gravity:      gravity,
		LifeLength:   life,
		Transparency: 1,
		Rotation:     rotAngle,
		Scale:        scale,
	}
}

type Particle struct {
	Position     [3]float32
	Velocity     [3]float32
	Gravity      float32
	LifeLength   float32
	Transparency float32
	Rotation     float32
	Scale        float32
	elapsedTime  float32

	tmp1, tmp2 mgl32.Mat4
}

func (p *Particle) Update(elapsed float32) bool {
	p.Velocity[1] += Gravity * p.Gravity * elapsed
	change := mgl32.Vec3(p.Velocity).Mul(elapsed)
	p.Position[0] += change[0]
	p.Position[1] += change[1]
	p.Position[2] += change[2]
	p.elapsedTime += elapsed
	p.Transparency = 1 - p.elapsedTime/p.LifeLength
	return p.elapsedTime < p.LifeLength
}

func (t *Particle) Transform(camera components.Viewable) mgl32.Mat4 {
	translateMatrix := mgl32.Translate3D(t.Position[0], t.Position[1], t.Position[2])

	view := camera.View()
	// always face camera
	translateMatrix.Set(0, 0, view.At(0, 0))
	translateMatrix.Set(1, 0, view.At(0, 1))
	translateMatrix.Set(2, 0, view.At(0, 2))
	translateMatrix.Set(0, 1, view.At(1, 0))
	translateMatrix.Set(1, 1, view.At(1, 1))
	translateMatrix.Set(2, 1, view.At(1, 2))
	translateMatrix.Set(0, 2, view.At(2, 0))
	translateMatrix.Set(1, 2, view.At(2, 1))
	translateMatrix.Set(2, 2, view.At(2, 2))

	// @todo optimise this
	//rotationMatrix := mgl32.HomogRotate3D(mgl32.DegToRad(t.Rotation), mgl32.Vec3{0, 0, 1})
	rotationMatrix := mgl32.Ident4()
	scaleMatrix := mgl32.Scale3D(t.Scale, t.Scale, t.Scale)

	math.Mul4(rotationMatrix, scaleMatrix, &t.tmp1)
	math.Mul4(translateMatrix, t.tmp1, &t.tmp2)
	math.Mul4(view, t.tmp2, &t.tmp1)
	math.Mul4(camera.Projection(), t.tmp1, &t.tmp2)
	return t.tmp2
}
