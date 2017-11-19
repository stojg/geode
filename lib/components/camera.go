package components

import "github.com/go-gl/mathgl/mgl32"

func NewCamera(projection mgl32.Mat4) *Camera {
	return &Camera{
		projection: projection,
	}
}

type Camera struct {
	GameComponent

	projection mgl32.Mat4
}

func (c *Camera) AddToEngine(engine Engine) {
	engine.GetRenderingEngine().AddCamera(c)
}

// https://github.com/BennyQBD/GDX/blob/772b4c75c29c65e36ce100755da0ea483c80cee8/GDX/camera.cpp
func (c *Camera) GetViewProjection() mgl32.Mat4 {
	cameraRotation := c.Transform().TransformedRot().Conjugate().Mat4()
	cameraPos := c.Transform().TransformedPos().Mul(-1)
	cameraTranslation := mgl32.Translate3D(cameraPos[0], cameraPos[1], cameraPos[2])
	return c.projection.Mul4(cameraRotation.Mul4(cameraTranslation))
}
