package components

import (
	"math"
	"time"

	"github.com/go-gl/mathgl/mgl32"
)

func NewCamera(fovy float32, width, height int, near, far float32) *Camera {

	return &Camera{
		projection: mgl32.Perspective(mgl32.DegToRad(fovy), float32(width/height), near, far),
	}
}

type Camera struct {
	GameComponent

	projection mgl32.Mat4
	planes     [6][4]float32
}

func (c *Camera) Planes() [6][4]float32 {
	return c.planes
}

func (c *Camera) Pos() mgl32.Vec3 {
	return c.parent.Transform().TransformedPos()
}

func (c *Camera) Rot() mgl32.Quat {
	return c.parent.Transform().TransformedRot()
}

func (c *Camera) AddToEngine(e Engine) {
	e.RenderingEngine().State().AddCamera(c)
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

// @todo this should probably be replaced with a function that returns the planes

func (c *Camera) Update(ts time.Duration) {
	c.planes = extractPlanesFromProjection(c.Projection().Mul4(c.View()), true)
	//not := IsVisible(c.planes, [3]float32{0, 1.144242, 0}, mgl32.Mat4{})
	//if !not {
	//	fmt.Println("not vis")
	//} else {
	//	fmt.Println("vis")
	//}
}

func extractPlanesFromProjection(projection mgl32.Mat4, normalise bool) [6][4]float32 {
	var res [6][4]float32
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
		normalisePlane(&res[0])
		normalisePlane(&res[1])
		normalisePlane(&res[2])
		normalisePlane(&res[3])
		normalisePlane(&res[4])
		normalisePlane(&res[5])
	}
	return res
}

func normalisePlane(a *[4]float32) {
	l := float32(math.Sqrt(float64(a[0]*a[0] + a[1]*a[1] + a[2]*a[2])))
	a[0] /= l
	a[1] /= l
	a[2] /= l
	a[3] /= l
}
