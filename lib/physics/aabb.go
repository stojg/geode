package physics

import "github.com/go-gl/mathgl/mgl32"

type AABB struct {
	c mgl32.Vec3 // center point
	r mgl32.Vec3 // radius, or half width extents (rx, ry, rz)
}

func (A *AABB) C() mgl32.Vec3 {
	return A.c
}

func (A *AABB) SetC(c mgl32.Vec3) {
	A.c = c
}

func (A *AABB) R() mgl32.Vec3 {
	return A.r
}

func (A *AABB) SetR(r mgl32.Vec3) {
	A.r = r
}
