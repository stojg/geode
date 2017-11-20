package core

import (
	"fmt"
	"time"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/stojg/graphics/lib/components"
	"github.com/stojg/graphics/lib/rendering"
)

type Drawable interface {
	Draw()
}

func Main(log Logger) error {
	width := 800
	height := 600

	engine, err := NewEngine(width, height, "graphics")
	if err != nil {
		return err
	}

	cameraObject := NewGameObject()
	cameraObject.Transform().SetPos(mgl32.Vec3{0, 0, 10})
	cameraObject.Transform().LookAt(mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, 1, 0})
	cameraObject.AddComponent(components.NewCamera(70, float32(width), float32(height), 0.01, 100))
	cameraObject.AddComponent(&components.FreeMove{})
	cameraObject.AddComponent(components.NewFreelook(float32(width), float32(height)))
	engine.Game().AddObject(cameraObject)

	cubeMesh, err := rendering.NewMesh("res/meshes/cube/model.obj")
	if err != nil {
		fmt.Printf("Model loading failed: %v", err)
		return err
	}

	whiteMaterial := rendering.NewMaterial()
	whiteMaterial.AddTexture("diffuse", rendering.NewTexture("res/textures/white.png"))

	floor := NewGameObject()
	floor.Transform().SetScale(mgl32.Vec3{100, 0.01, 100})
	floor.Transform().SetPos(mgl32.Vec3{0, -3, 0})
	floor.AddComponent(components.NewMeshRenderer(cubeMesh, whiteMaterial))
	engine.Game().AddObject(floor)

	{
		light := NewGameObject()
		light.Transform().SetPos(mgl32.Vec3{-3, 3, 2})
		light.Transform().SetScale(mgl32.Vec3{0.2, 0.2, 0.2})
		light.AddComponent(components.NewMeshRenderer(cubeMesh, whiteMaterial))

		dirLight := components.NewBaseLight(mgl32.Vec3{0.5, 0.9, 1}, 1)
		dirLight.SetShader(rendering.NewShader("forward_point"))
		light.AddComponent(dirLight)
		engine.Game().AddObject(light)
	}

	{
		light := NewGameObject()
		light.Transform().SetPos(mgl32.Vec3{3, 3, 2})
		light.Transform().SetScale(mgl32.Vec3{0.2, 0.2, 0.2})
		light.AddComponent(components.NewMeshRenderer(cubeMesh, whiteMaterial))

		dirLight := components.NewBaseLight(mgl32.Vec3{1, 0.9, 0.5}, 1)
		dirLight.SetShader(rendering.NewShader("forward_point"))
		light.AddComponent(dirLight)
		engine.Game().AddObject(light)
	}

	material := rendering.NewMaterial()
	material.AddTexture("diffuse", rendering.NewTexture("res/textures/test.png"))

	cubeRenderer := components.NewMeshRenderer(cubeMesh, material)
	cube := NewGameObject()
	cube.AddComponent(cubeRenderer)
	cube.AddComponent(&components.Rotator{})
	cube.Transform().SetPos(mgl32.Vec3{0, 0, 0})
	engine.Game().AddObject(cube)

	engine.Start()
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
