package components

import "github.com/go-gl/mathgl/mgl32"

func up() mgl32.Vec3 {
	return mgl32.Vec3{0, 1, 0}
}

func down() mgl32.Vec3 {
	return mgl32.Vec3{0, -1, 0}
}

func right() mgl32.Vec3 {
	return mgl32.Vec3{1, 0, 0}
}

func left() mgl32.Vec3 {
	return mgl32.Vec3{-1, 0, 0}
}

func forward() mgl32.Vec3 {
	return mgl32.Vec3{0, 0, 1}
}

func back() mgl32.Vec3 {
	return mgl32.Vec3{0, 0, -1}
}

func vec3(x, y, z float32) mgl32.Vec3 {
	return mgl32.Vec3{x, y, z}
}
