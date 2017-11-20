package core

import (
	"fmt"
	"time"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/stojg/graphics/lib/components"
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
	cameraObject.Transform().SetPos(mgl32.Vec3{0, 2, 10})
	cameraObject.Transform().LookAt(mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, 1, 0})
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
		light := NewGameObject()
		light.Transform().SetPos(mgl32.Vec3{-3, 3, 2})
		light.Transform().SetScale(mgl32.Vec3{0.1, 0.1, 0.1})
		light.AddComponent(components.NewRotator(mgl32.Vec3{1, 1, 1}, 90))
		dirLight := components.NewBaseLight(mgl32.Vec3{0.4, 0.9, 1}, 1)
		dirLight.SetShader(rendering.NewShader("forward_point"))
		light.AddComponent(dirLight)
		LoadModel(light, "res/meshes/cube/model.obj", whiteMaterial)
		engine.Game().AddObject(light)
	}

	{
		light := NewGameObject()
		light.Transform().SetPos(mgl32.Vec3{3, 3, 2})
		light.Transform().SetScale(mgl32.Vec3{0.1, 0.1, 0.1})
		light.AddComponent(components.NewRotator(mgl32.Vec3{1, 1, 1}, 90))

		dirLight := components.NewBaseLight(mgl32.Vec3{1, 0.9, 0.4}, 1)
		dirLight.SetShader(rendering.NewShader("forward_point"))
		light.AddComponent(dirLight)
		LoadModel(light, "res/meshes/cube/model.obj", whiteMaterial)
		engine.Game().AddObject(light)
	}

	bot := NewGameObject()
	bot.Transform().SetPos(mgl32.Vec3{0, 0, 0})
	bot.AddComponent(components.NewRotator(mgl32.Vec3{0, 1, 0}, 23))
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
	for _, data := range objData {
		mesh := rendering.NewMesh()
		mesh.SetVertices(rendering.ConvertToVertices(data))
		g.AddComponent(components.NewMeshRenderer(mesh, material))
	}
	return nil
}

func NewGame() *Game {
	g := &Game{}
	return g
}

type Game struct {
	root *GameObject
}

func (g *Game) SetEngine(engine *Engine) {
	g.RootObject().SetEngine(engine)
}

func (g *Game) AddObject(object *GameObject) {
	g.RootObject().AddChild(object)
}

func (g *Game) Input(elapsed time.Duration) {
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
