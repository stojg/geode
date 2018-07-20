package physics

// https://github.com/BennyQBD/3DGameEngine/blob/225fa8baf6637756ba03ccbc0444bf7751d87dbb/src/com/base/engine/core/Transform.java

import (
	"math"

	x "github.com/stojg/graphics/lib/math"

	"github.com/go-gl/mathgl/mgl32"
)

func NewTransform() *Transform {
	return &Transform{
		pos:   mgl32.Vec3{0, 0, 0},
		rot:   mgl32.QuatIdent(),
		scale: mgl32.Vec3{1, 1, 1},

		oldPos:       mgl32.Vec3{math.MaxFloat32 - 1, math.MaxFloat32 - 1, math.MaxFloat32 - 1},
		oldRot:       mgl32.Quat{W: math.MaxFloat32 - 1, V: mgl32.Vec3{math.MaxFloat32 - 1, math.MaxFloat32 - 1, math.MaxFloat32 - 1}},
		parentMatrix: mgl32.Ident4(),
	}
}

type Transform struct {
	parent       *Transform
	parentMatrix mgl32.Mat4

	pos   mgl32.Vec3
	rot   mgl32.Quat
	scale mgl32.Vec3

	oldPos   mgl32.Vec3
	oldRot   mgl32.Quat
	oldScale mgl32.Vec3

	transformation mgl32.Mat4
}

func (t *Transform) Update() {
	t.oldPos = t.pos
	t.oldRot = t.rot
	t.oldScale = t.scale
	t.transformation = t.calcTransformation()
}

// angle should be in radians
func (t *Transform) Rotate(axis mgl32.Vec3, angle float32) {
	t.rot = t.rot.Mul(mgl32.QuatRotate(angle, axis)).Normalize()
}

func (t *Transform) MoveBy(add mgl32.Vec3) {
	t.pos = t.pos.Add(add)
}

func (t *Transform) LookAt(point mgl32.Vec3, up mgl32.Vec3) {
	t.rot = t.LookAtRotation(point, up)
}

func (t *Transform) LookAtRotation(point mgl32.Vec3, up mgl32.Vec3) mgl32.Quat {
	eye := t.pos
	center := point

	direction := center.Sub(eye).Normalize()

	// Find the rotation between the front of the object (that we assume towards Z-,
	// but this depends on your model) and the desired direction
	rotDir := mgl32.QuatBetweenVectors(mgl32.Vec3{0, 0, -1}, direction)

	// Recompute up so that it's perpendicular to the direction
	// You can skip that part if you really want to force up
	right := direction.Cross(up)
	up = right.Cross(direction)

	// Because of the 1rst rotation, the up is probably completely screwed up.
	// Find the rotation between the "up" of the rotated object, and the desired up
	upCur := rotDir.Rotate(mgl32.Vec3{0, 1, 0})
	rotUp := mgl32.QuatBetweenVectors(upCur, up)

	rotTarget := rotUp.Mul(rotDir) // remember, in reverse order.

	return rotTarget // camera rotation should be inversed!
}

func (t *Transform) HasChanged() bool {
	if t.parent != nil && t.parent.HasChanged() {
		return true
	}

	if !t.pos.ApproxEqual(t.oldPos) {
		return true
	}

	if !t.rot.ApproxEqual(t.rot) {
		return true
	}

	if !t.scale.ApproxEqual(t.scale) {
		return true
	}
	return false
}

var tmp1, tmp2 mgl32.Mat4

func (t *Transform) calcTransformation() mgl32.Mat4 {
	translationMatrix := mgl32.Translate3D(t.pos[0], t.pos[1], t.pos[2])
	rotationMatrix := t.rot.Mat4()
	scaleMatrix := mgl32.Scale3D(t.scale[0], t.scale[1], t.scale[2])

	x.Mul4(rotationMatrix, scaleMatrix, &tmp1)
	x.Mul4(translationMatrix, tmp1, &tmp2)
	x.Mul4(t.ParentMatrix(), tmp2, &tmp1)
	return tmp1
	//return t.ParentMatrix().Mul4(translationMatrix.Mul4(rotationMatrix.Mul4(scaleMatrix)))
}

func (t *Transform) Transformation() mgl32.Mat4 {
	return t.transformation
}

func (t *Transform) ParentMatrix() mgl32.Mat4 {
	if t.parent != nil {
		t.parentMatrix = t.parent.Transformation()
	}
	return t.parentMatrix
}

func (t *Transform) SetParent(parent *Transform) {
	t.parent = parent
}

func (t *Transform) TransformedPos() mgl32.Vec3 {
	return mgl32.TransformCoordinate(t.pos, t.ParentMatrix())
}

func (t *Transform) TransformedRot() mgl32.Quat {
	parentRot := mgl32.QuatIdent()
	if t.parent != nil {
		parentRot = t.parent.TransformedRot()
	}
	return parentRot.Mul(t.rot)
}

func (t *Transform) Pos() mgl32.Vec3 {
	return t.pos
}

func (t *Transform) SetPos(pos mgl32.Vec3) {
	t.pos = pos
}

func (t *Transform) Rot() mgl32.Quat {
	return t.rot
}

func (t *Transform) SetRot(rot mgl32.Quat) {
	t.rot = rot
}

func (t *Transform) Scale() mgl32.Vec3 {
	return t.scale
}

func (t *Transform) SetScale(scale mgl32.Vec3) {
	t.scale = scale
}
