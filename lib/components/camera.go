package components

import (
	"time"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/stojg/geode/lib/geometry"
)

func NewCamera(fovy float32, width, height int, near, far float32) *Camera {
	return &Camera{
		projection: mgl32.Perspective(mgl32.DegToRad(fovy), float32(width/height), near, far),
		view:       mgl32.Ident4(),
	}
}

type Camera struct {
	GameComponent

	projection mgl32.Mat4
	view       mgl32.Mat4
	planes     geometry.Frustum
}

func (c *Camera) Frustum() geometry.Frustum {
	return c.planes
}

func (c *Camera) Pos() mgl32.Vec3 {
	return c.parent.Transform().TransformedPos()
}

func (c *Camera) Rot() mgl32.Quat {
	return c.parent.Transform().TransformedRot()
}

func (c *Camera) AddToEngine(e RenderState) {
	e.SetCamera(c)
}

func (c *Camera) Projection() mgl32.Mat4 {
	return c.projection
}

func (c *Camera) View() mgl32.Mat4 {
	return c.view
}

// @todo this should probably be replaced with a function that returns the planes

func (c *Camera) Update(ts time.Duration) {
	c.view = calcView(c)
	c.planes = extractPlanesFromProjection(c.Projection().Mul4(c.View()), true)
}

func calcView(c *Camera) mgl32.Mat4 {
	//This comes from the conjugate rotation because the world should appear to rotate opposite to the camera's rotation.
	cameraRotation := c.Transform().TransformedRot().Conjugate().Mat4()
	//Similarly, the translation is inverted because the world appears to move opposite to the camera's movement.
	cameraPos := c.Transform().TransformedPos().Mul(-1)
	cameraTranslation := mgl32.Translate3D(cameraPos[0], cameraPos[1], cameraPos[2])
	return cameraRotation.Mul4(cameraTranslation)
}

func extractPlanesFromProjection(projection mgl32.Mat4, normalise bool) geometry.Frustum {
	var res geometry.Frustum
	for i := 0; i < 4; i++ {
		f := projection[3+i*4]
		res[0][i] = f + projection[i*4+0] // left
		res[1][i] = f - projection[i*4+0] // right
		res[2][i] = f + projection[i*4+1] // bottom
		res[3][i] = f - projection[i*4+1] // top
		res[4][i] = f + projection[i*4+2] // near
		res[5][i] = f - projection[i*4+2] // far
	}

	if normalise {
		res.Normalise()
	}
	return res
}
