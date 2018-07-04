package core

import (
	"github.com/go-gl/mathgl/mgl32"
)

func IsVisible(planes [6][4]float32, aabb [3][2]float32, transform mgl32.Mat4) bool {
	return aabbInFrustum(planes, aabb, transform)
}

func aabbInFrustum(planes [6][4]float32, aabb [3][2]float32, transform mgl32.Mat4) bool {

	var d float32
	min := transform.Mul4x1(mgl32.Vec4{-aabb[0][0] + aabb[0][1], -aabb[1][0] + aabb[1][1], -aabb[2][0] + aabb[2][1], 1})
	max := transform.Mul4x1(mgl32.Vec4{aabb[0][0] + aabb[0][1], aabb[1][0] + aabb[1][1], aabb[2][0] + aabb[2][1], 1})

	// check box outside/inside of frustum
	for i := 0; i < 6; i++ {
		out := 0
		dotPlane(planes[i], [3]float32{min[0], min[1], min[2]}, &d)
		if d < 0.0 {
			out++
		}
		dotPlane(planes[i], [3]float32{max[0], min[1], min[2]}, &d)
		if d < 0.0 {
			out++
		}
		dotPlane(planes[i], [3]float32{min[0], max[1], min[2]}, &d)
		if d < 0.0 {
			out++
		}
		dotPlane(planes[i], [3]float32{max[0], max[1], min[2]}, &d)
		if d < 0.0 {
			out++
		}
		dotPlane(planes[i], [3]float32{min[0], min[1], max[2]}, &d)
		if d < 0.0 {
			out++
		}
		dotPlane(planes[i], [3]float32{max[0], min[1], max[2]}, &d)
		if d < 0.0 {
			out++
		}
		dotPlane(planes[i], [3]float32{min[0], max[1], max[2]}, &d)
		if d < 0.0 {
			out++
		}
		dotPlane(planes[i], [3]float32{max[0], max[1], max[2]}, &d)
		if d < 0.0 {
			out++
		}
		if out == 8 {
			return false
		}
	}

	return true
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
