package core

import (
	"fmt"
	"time"

	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/stojg/graphics/lib/components"
	"github.com/stojg/graphics/lib/input"
	"github.com/stojg/graphics/lib/lights"
	"github.com/stojg/graphics/lib/rendering"
	"github.com/stojg/graphics/lib/rendering/loader"
)

func Main(log Logger) error {
	width := 800
	height := 600

	engine, err := NewEngine(width, height, "graphics")
	if err != nil {
		return err
	}

	cameraObject := NewGameObject()
	cameraObject.Transform().SetPos(mgl32.Vec3{8, 5, 8})
	cameraObject.Transform().LookAt(mgl32.Vec3{0, 1, 0}, mgl32.Vec3{0, 1, 0})
	cameraObject.AddComponent(components.NewCamera(70, float32(width), float32(height), 0.01, 100))
	cameraObject.AddComponent(&components.FreeMove{})
	cameraObject.AddComponent(components.NewFreelook(float32(width), float32(height)))
	engine.Game().AddObject(cameraObject)

	whiteMaterial := rendering.NewMaterial()
	whiteMaterial.AddTexture("diffuse", rendering.NewTexture("res/textures/white.png"))

	floor := NewGameObject()
	floor.Transform().SetScale(mgl32.Vec3{100, 0.01, 100})
	floor.Transform().SetPos(mgl32.Vec3{0, -0.005, 0})
	LoadModel(floor, "res/meshes/cube/model.obj", whiteMaterial)
	engine.Game().AddObject(floor)

	{
		dirLight := NewGameObject()
		dirLight.Transform().SetPos(mgl32.Vec3{2, 6, -1})
		dirLight.AddComponent(lights.NewDirectional(0.99, 0.98, 0.7, 1))
		engine.Game().AddObject(dirLight)
	}

	bot := NewGameObject()
	bot.Transform().SetPos(mgl32.Vec3{0, 0, 0})
	bot.AddComponent(components.NewRotator(mgl32.Vec3{0, -1, 0}, 23))
	LoadModel(bot, "res/meshes/sphere_bot/model.obj", whiteMaterial)
	engine.Game().AddObject(bot)

	engine.Start()
	return nil
}

func LoadModel(g *GameObject, obj string, material *rendering.Material) error {
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
