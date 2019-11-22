package core

import (
	"fmt"
	"time"

	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/stojg/geode/lib/components"
	"github.com/stojg/geode/lib/input"
	"github.com/stojg/geode/lib/rendering"
)

func NewScene() *Scene {
	g := &Scene{
		vsync:    false,
		exposure: 4.0,
		effects:  true,
	}
	return g
}

type Scene struct {
	root     *GameObject
	vsync    bool
	effects  bool
	exposure float32
	state    components.RenderState
}

func (g *Scene) SetState(state components.RenderState) {
	g.state = state
	g.state.SetInteger("effects", 1)
	g.state.SetFloat("x_exposure", g.exposure)
	g.rootObject().SetState(state)
}

func (g *Scene) AddObject(object components.Object) {
	g.rootObject().AddChild(object)
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

	if input.KeyDown(glfw.KeyLeftBracket) {
		g.exposure -= 0.2
		g.state.SetFloat("x_exposure", g.exposure)
		fmt.Printf("setting exposure to %0.1f\n", g.exposure)
	}
	if input.KeyDown(glfw.KeyRightBracket) {
		g.exposure += 0.2
		g.state.SetFloat("x_exposure", g.exposure)
		fmt.Printf("setting exposure to %0.1f\n", g.exposure)
	}
}

func (g *Scene) Update(elapsed time.Duration) {
	g.rootObject().Update(elapsed)
}

func (g *Scene) Render(r *rendering.Renderer) {
	r.Render(g.rootObject())
}

func (g *Scene) rootObject() *GameObject {
	if g.root == nil {
		g.root = NewGameObject(components.ResourceNA)
	}
	return g.root
}
