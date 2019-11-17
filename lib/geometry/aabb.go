package geometry

import "github.com/go-gl/mathgl/mgl32"

type AABB struct {
	c mgl32.Vec3 // center point
	r mgl32.Vec3 // radius, or half width extents (rx, ry, rz)
}

// C returns the center point of the AABB
func (a *AABB) C() mgl32.Vec3 {
	return a.c
}

// SetC sets the center point of the AABB
func (a *AABB) SetC(c mgl32.Vec3) {
	a.c = c
}

// R return the radius, or half width extents, of the AABB
func (a *AABB) R() mgl32.Vec3 {
	return a.r
}

// SetR set set the radius, or half width extents, of the AABB
func (a *AABB) SetR(r mgl32.Vec3) {
	a.r = r
}
