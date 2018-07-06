package components

import (
	"time"
)

func NewHeadHeight(t Terrain) *HeadHeight {
	return &HeadHeight{
		terrain: t,
	}
}

type HeadHeight struct {
	GameComponent
	terrain Terrain
}

func (c *HeadHeight) Update(elapsed time.Duration) {
	temp := c.Transform().Pos()
	x, z := temp[0], temp[2]
	temp[1] = 1.6 + c.terrain.Height(x, z)
	c.Transform().SetPos(temp)
}
