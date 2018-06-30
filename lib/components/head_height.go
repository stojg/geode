package components

import (
	"time"
)

type HeadHeight struct {
	GameComponent
	Terrain Terrain
}

func (c *HeadHeight) Update(elapsed time.Duration) {
	temp := c.Transform().Pos()
	x, z := temp[0], temp[2]
	temp[1] = 1.6 + c.Terrain.Height(x, z)
	c.Transform().SetPos(temp)
}
