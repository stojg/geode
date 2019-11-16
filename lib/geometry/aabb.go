package geometry

import "github.com/go-gl/mathgl/mgl32"

type AABB struct {
	c mgl32.Vec3 // center point
	r mgl32.Vec3 // radius, or half width extents (rx, ry, rz)
}

// C returns the center point of the AABB
func (A *AABB) C() mgl32.Vec3 {
	return A.c
}

// SetC sets the center point of the AABB
func (A *AABB) SetC(c mgl32.Vec3) {
	A.c = c
}

// R return the radius, or half width extents, of the AABB
func (A *AABB) R() mgl32.Vec3 {
	return A.r
}

// SetR set set the radius, or half width extents, of the AABB
func (A *AABB) SetR(r mgl32.Vec3) {
	A.r = r
}
