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
	width := 400
	height := 300

	engine, err := NewEngine(width, height, "graphics")
	if err != nil {
		return err
	}

	cameraObject := NewGameObject()
	cameraObject.AddComponent(components.NewCamera(70, float32(width), float32(height), 0.01, 100))
	cameraObject.AddComponent(&components.FreeMove{})
	cameraObject.AddComponent(components.NewFreelook(float32(width), float32(height)))
	cameraObject.Transform().SetPos(mgl32.Vec3{0, 0, 2})
	cameraObject.Transform().LookAt(mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, 1, 0})
	engine.Game().AddObject(cameraObject)

	cubeMesh, err := rendering.NewMesh("res/meshes/cube/model.obj")
	if err != nil {
		fmt.Printf("Model loading failed: %v", err)
		return err
	}
	material := rendering.NewMaterial()
	material.AddTexture("diffuse", rendering.NewTexture("res/textures/test.png"))

	meshRenderer := components.NewMeshRenderer(cubeMesh, material)
	cube := NewGameObject()
	cube.AddComponent(meshRenderer)
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
