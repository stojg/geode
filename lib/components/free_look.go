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

	if delta[0] == 0 && delta[1] == 0 {
		return
	}

	const sensitivity float32 = 0.01
	c.yaw -= delta[0] * sensitivity
	c.pitch += delta[1] * sensitivity

	kQuat := mgl32.Quat{W: 1, V: mgl32.Vec3{c.pitch, c.yaw, c.roll}}
	camQUat := kQuat.Mul(c.Transform().Rot())
	camQUat.Normalize()
	c.Transform().SetRot(camQUat)

	c.pitch, c.yaw, c.roll = 0, 0, 0
	input.SetCursorPosition(centerPosition[0], centerPosition[1])
}
