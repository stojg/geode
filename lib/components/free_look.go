package components

import (
	"time"

	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/stojg/graphics/lib/input"
)

type FreeLook struct {
	GameComponent

	yaw, pitch, roll float32
	locked           bool
}

func (c *FreeLook) Update(elapsed time.Duration) {

}

func (c *FreeLook) Input(elapsed time.Duration) {

	centerPosition := mgl32.Vec2{800 / 2, 600 / 2}
	if input.ButtonDown(glfw.MouseButton1) {
		c.locked = true
		input.HideCursor()
		input.SetCursorPosition(centerPosition[0], centerPosition[1])
	}

	if input.KeyDown(glfw.KeySpace) {
		c.locked = false
		input.SetCursorPosition(centerPosition[0], centerPosition[1])
		input.ShowCursor()
	}

	if !c.locked {
		return
	}

	delta := mgl32.Vec2(input.CursorPosition()).Sub(centerPosition)
	if delta.Len() == 0 {
		return
	}
	input.SetCursorPosition(centerPosition[0], centerPosition[1])

	const sensitivity float32 = 0.5

	yaw := mgl32.DegToRad(-delta[0]) * sensitivity
	pitch := mgl32.DegToRad(delta[1]) * sensitivity
	var roll float32 = 0.0

	//temp := mgl32.QuatRotate(1, mgl32.Vec3{pitch, yaw, roll})
	temp := mgl32.Quat{W: 10, V: mgl32.Vec3{pitch, yaw, roll}}.Normalize()
	camQUat := temp.Mul(c.Transform().Rot()).Normalize()
	c.Transform().SetRot(camQUat)

}
