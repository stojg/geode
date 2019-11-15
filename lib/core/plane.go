package core

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/stojg/geode/lib/components"
)

func IsVisible(planes [6][4]float32, aabb components.AABB, transform mgl32.Mat4) bool {
	return aabbInFrustum(planes, aabb, transform)
}

func aabbInFrustum(planes [6][4]float32, aabb components.AABB, transform mgl32.Mat4) bool {

	for i := 0; i < 6; i++ {

		// compute distance of box center from plane
		d := dot(planes[i], aabb.C())

		// compute the projection interval radius of b onto plane
		r := aabb.R()[0]*abs(planes[i][0]) + aabb.R()[1]*abs(planes[i][1]) + aabb.R()[2]*abs(planes[i][2])

		if d+r <= -planes[i][3] {
			return false
		}
	}
	return true
}

func dot(a [4]float32, b mgl32.Vec3) float32 {
	x := a[0] * b[0]
	y := a[1] * b[1]
	z := a[2] * b[2]
	return x + y + z
}

func abs(x float32) float32 {
	if x > 0 {
		return x
	}
	return -x
}

func pointInFrustum(planes [6][4]float32, pt [3]float32) bool {
	var d float32
	for i := 0; i < 6; i++ {
		dotPlane(planes[i], pt, &d)
		if d <= 0 {
			return false
		}
	}
	return true
}

func sphereInFrustum(planes [6][4]float32, pt [3]float32, radius float32) bool {
	var d float32
	for i := 0; i < 6; i++ {
		dotPlane(planes[i], pt, &d)
		if d <= -radius {
			return false
		}
	}
	return true
}

func dotPlane(plane [4]float32, pt [3]float32, result *float32) {
	a := plane[0] * pt[0]
	b := plane[1] * pt[1]
	c := plane[2] * pt[2]
	*result = a + b + c + plane[3]
}
