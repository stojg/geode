package core

import (
	"fmt"
	"time"

	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/stojg/graphics/lib/components"
	"github.com/stojg/graphics/lib/input"
	"github.com/stojg/graphics/lib/rendering"
	"github.com/stojg/graphics/lib/rendering/loader"
)

func LoadModel(g *GameObject, obj string, material components.Material) error {
	objData, err := loader.Load(obj)
	if err != nil {
		fmt.Printf("Model loading failed: %v", err)
		return err
	}
	for i, data := range objData {
		mesh := rendering.NewMesh()
		mesh.SetVertices(rendering.ConvertToVertices(data))
		fmt.Printf("LoadModel: %s.%d has %d vertices\n", obj, i, mesh.NumVertices())
		g.AddComponent(components.NewMeshRenderer(mesh, material))
	}
	return nil
}

func NewGame() *Game {
	g := &Game{
		vsync: true,
	}
	return g
}

type Game struct {
	root  *GameObject
	vsync bool
}

func (g *Game) SetEngine(engine *Engine) {
	g.RootObject().SetEngine(engine)
}

func (g *Game) AddObject(object *GameObject) {
	g.RootObject().AddChild(object)
}

func (g *Game) Input(elapsed time.Duration) {
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
	g.RootObject().InputAll(elapsed)
}

func (g *Game) Update(elapsed time.Duration) {
	g.RootObject().UpdateAll(elapsed)
}

func (g *Game) Render(r *rendering.Engine) {
	r.Render(g.RootObject())
}

func (g *Game) RootObject() *GameObject {
	if g.root == nil {
		g.root = NewGameObject()
	}
	return g.root
}
