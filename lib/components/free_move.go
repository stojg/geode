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
	const speed float32 = 10
	amount := float32(elapsed.Seconds()) * speed

	if input.Key(glfw.KeyW) {
		c.Move(c.Transform().Rot().Rotate(back()), amount)
	} else if input.Key(glfw.KeyS) {
		c.Move(c.Transform().Rot().Rotate(forward()), amount)
	}

	if input.Key(glfw.KeyA) {
		c.Move(c.Transform().Rot().Rotate(left()), amount)
	} else if input.Key(glfw.KeyD) {
		c.Move(c.Transform().Rot().Rotate(right()), amount)
	}
}

func (c *FreeMove) Move(dir mgl32.Vec3, amount float32) {
	c.Transform().MoveBy(dir.Mul(amount))
}
