package core

import (
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

	mesh := rendering.NewMesh()
	vertices := []rendering.Vertex{
		{Pos: [3]float32{-0.5, -0.5, 0.0}},
		{Pos: [3]float32{0.5, -0.5, 0.0}},
		{Pos: [3]float32{0, 0.5, 0.0}},
	}
	mesh.AddVertices(vertices)

	meshRenderer := components.NewMeshRenderer(mesh)
	triangleObject := NewGameObject()
	triangleObject.AddComponent(meshRenderer)
	triangleObject.AddComponent(&components.Rotator{})
	triangleObject.Transform().SetPos(mgl32.Vec3{0, 0, 0})
	engine.Game().AddObject(triangleObject)

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
