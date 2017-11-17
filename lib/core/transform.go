package core
// https://github.com/BennyQBD/3DGameEngine/blob/225fa8baf6637756ba03ccbc0444bf7751d87dbb/src/com/base/engine/core/Transform.java

import (

	"github.com/go-gl/mathgl/mgl32"
)


func NewTransform() *Transform {
	return &Transform{
		pos:   mgl32.Vec3{0, 0, 0},
		rot:   mgl32.Quat{W: 1, V: mgl32.Vec3{0, 0, 0}},
		scale: mgl32.Vec3{0, 0, 0},
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
}

func (t *Transform) Rotate(axis mgl32.Vec3, angle float32) {
	t.rot = mgl32.Quat{W: angle, V: axis}.Mul(t.rot).Normalize()
}

func (t *Transform) LookAt(point mgl32.Vec3, up mgl32.Vec3) {
	t.rot = t.GetLookAtRotation(point, up)
}

func (t *Transform) GetLookAtRotation(point mgl32.Vec3, up mgl32.Vec3) mgl32.Quat {
	return mgl32.QuatLookAtV(t.pos, point, up)
}

func (t *Transform) HasChanged() bool {
	if t.parent != nil && t.parent.HasChanged() {
		return true
	}

	if t.pos.ApproxEqual(t.oldPos) {
		return true
	}

	if t.rot.ApproxEqual(t.rot) {
		return true
	}

	if t.scale.ApproxEqual(t.scale) {
		return true
	}
	return false
}

func (t *Transform) GetTransformation() mgl32.Mat4 {
	/*
	Matrix4f translationMatrix = new Matrix4f().InitTranslation(m_pos.GetX(), m_pos.GetY(), m_pos.GetZ());
		Matrix4f rotationMatrix = m_rot.ToRotationMatrix();
		Matrix4f scaleMatrix = new Matrix4f().InitScale(m_scale.GetX(), m_scale.GetY(), m_scale.GetZ());

		return GetParentMatrix().Mul(translationMatrix.Mul(rotationMatrix.Mul(scaleMatrix)));
	 */
	 return mgl32.Mat4{}
}

func (t *Transform) GetTransformedPos() {
	//return t.ParentMatrix().Transform(t.pos)
}

func (t *Transform) GetTransformedRot() mgl32.Quat {
	parentRot := mgl32.QuatIdent()
	if t.parent != nil {
		parentRot = t.parent.GetTransformedRot()
	}
	return parentRot.Mul(t.rot)
}

func (t *Transform) ParentMatrix() mgl32.Mat4 {
	return t.parentMatrix
}

func (t *Transform) Scale() mgl32.Vec3 {
	return t.scale
}

func (t *Transform) SetScale(scale mgl32.Vec3) {
	t.scale = scale
}

func (t *Transform) Rot() mgl32.Quat {
	return t.rot
}

func (t *Transform) SetRot(rot mgl32.Quat) {
	t.rot = rot
}

func (t *Transform) Pos() mgl32.Vec3 {
	return t.pos
}

func (t *Transform) SetPos(pos mgl32.Vec3) {
	t.pos = pos
}

func (t *Transform) Update() {
	// @todo check if this is the first call to Update
	t.oldPos = t.pos
	t.oldRot = t.rot
	t.oldScale = t.scale
}
