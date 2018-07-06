package core

import (
	"fmt"
	"time"

	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/stojg/graphics/lib/components"
	"github.com/stojg/graphics/lib/input"
	"github.com/stojg/graphics/lib/rendering"
)

func NewScene() *Scene {
	g := &Scene{
		vsync: false,
	}
	return g
}

type Scene struct {
	root    *GameObject
	terrain *GameObject
	vsync   bool
	effects bool
	state   components.RenderState
}

func (g *Scene) SetState(state components.RenderState) {
	g.state = state
	g.state.SetInteger("effects", 0)
	g.rootObject().SetState(state)
}

func (g *Scene) AddObject(object *GameObject) {
	g.rootObject().AddChild(object)
}

func (g *Scene) AddTerrain(object *GameObject) {
	g.rootTerrain().AddChild(object)
}

func (g *Scene) Input(elapsed time.Duration) {
	if input.KeyDown(glfw.KeyRightShift) {
		g.vsync = !g.vsync
		if g.vsync {
			glfw.SwapInterval(1)
			fmt.Println("vsync on")
		} else {
			glfw.SwapInterval(0)
			fmt.Println("vsync off")
		}
	}

	if input.KeyDown(glfw.KeyF) {
		g.effects = !g.effects
		if g.effects {
			g.state.SetInteger("effects", 1)
			fmt.Println("effects on")
		} else {
			g.state.SetInteger("effects", 0)
			fmt.Println("effects off")
		}
	}

	g.rootObject().InputAll(elapsed)
}

func (g *Scene) Update(elapsed time.Duration) {
	g.rootObject().UpdateAll(elapsed)
	g.rootTerrain().UpdateAll(elapsed)
}

func (g *Scene) Render(r *rendering.Renderer) {
	r.Render(g.rootObject(), g.rootTerrain())
}

func (g *Scene) rootObject() *GameObject {
	if g.root == nil {
		g.root = NewGameObject()
	}
	return g.root
}

func (g *Scene) rootTerrain() *GameObject {
	if g.terrain == nil {
		g.terrain = NewGameObject()
	}
	return g.terrain
}
