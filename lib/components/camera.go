package components

import "github.com/go-gl/mathgl/mgl32"

func NewCamera(fov float32, width, height int, near, far float32) *Camera {

	return &Camera{
		projection: mgl32.Perspective(mgl32.DegToRad(fov), float32(width/height), near, far),
	}
}

type Camera struct {
	GameComponent

	projection mgl32.Mat4
}

func (c *Camera) Pos() mgl32.Vec3 {
	return c.parent.Transform().TransformedPos()
}

func (c *Camera) Rot() mgl32.Quat {
	return c.parent.Transform().TransformedRot()
}

func (c *Camera) AddToEngine(engine Engine) {
	engine.RenderingEngine().AddCamera(c)
}

func (c *Camera) Projection() mgl32.Mat4 {
	return c.projection
}

func (c *Camera) View() mgl32.Mat4 {
	//This comes from the conjugate rotation because the world should appear to rotate opposite to the camera's rotation.
	cameraRotation := c.Transform().TransformedRot().Conjugate().Mat4()
	//Similarly, the translation is inverted because the world appears to move opposite to the camera's movement.
	cameraPos := c.Transform().TransformedPos().Mul(-1)
	cameraTranslation := mgl32.Translate3D(cameraPos[0], cameraPos[1], cameraPos[2])
	return cameraRotation.Mul4(cameraTranslation)
}
