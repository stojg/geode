package components

import "github.com/go-gl/mathgl/mgl32"

func NewCamera(projection mgl32.Mat4) *Camera {
	return &Camera{
		projection: projection,
	}
}

type Camera struct {
	BaseComponent

	projection mgl32.Mat4
}

func (c *Camera) GetViewProjection() mgl32.Mat4 {
	//cameraRotation = c.Transform().TransformedRot().Conjugate().ToRotationMatrix()
	//cameraPos = c.Transform().TransformedPos().Mul(-1)
	// var cameraTranslation mgl32.Mat4 = mgl32.Mat4{}.InitTranslation(cameraPos.X(), cameraPos.Y(), cameraPos.Z())
	// return c.projection.Mul(cameraRotation.Mul(cameraTranslation))
	return mgl32.Mat4{}
}
