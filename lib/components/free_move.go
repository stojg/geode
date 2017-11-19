package components

import (
	"time"

	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/stojg/graphics/lib/input"
)

type FreeMove struct {
	GameComponent
}

func (c *FreeMove) Input(elapsed time.Duration) {
	const speed float32 = 1
	amount := float32(elapsed.Seconds()) * speed

	if input.Key(glfw.KeyW) {
		c.Move(c.Transform().Rot().Rotate(mgl32.Vec3{0, 0, -1}), amount)
	} else if input.Key(glfw.KeyS) {
		c.Move(c.Transform().Rot().Rotate(mgl32.Vec3{0, 0, 1}), amount)
	}

	if input.Key(glfw.KeyA) {
		c.Move(c.Transform().Rot().Rotate(mgl32.Vec3{-1, 0, 0}), amount)
	} else if input.Key(glfw.KeyD) {
		c.Move(c.Transform().Rot().Rotate(mgl32.Vec3{1, 0, 0}), amount)
	}
}

func (c *FreeMove) Move(dir mgl32.Vec3, amount float32) {
	c.Transform().SetPos(c.Transform().Pos().Add(dir.Mul(amount)))
}
