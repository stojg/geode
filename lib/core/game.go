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
	width := 800
	height := 600

	engine, err := NewEngine(width, height, "graphics")
	if err != nil {
		return err
	}

	projection := mgl32.Perspective(mgl32.DegToRad(70), float32(width/height), 0.01, 1000.0)
	cameraObject := NewGameObject()
	cameraObject.AddComponent(components.NewCamera(projection))
	cameraObject.AddComponent(&components.FreeMove{})
	cameraObject.AddComponent(&components.FreeLook{})
	cameraObject.Transform().SetPos(mgl32.Vec3{0, 0, 2})
	cameraObject.Transform().LookAt(mgl32.Vec3{0, 0, -2}, mgl32.Vec3{0, 1, 0})
	engine.Game().AddObject(cameraObject)

	mesh := rendering.NewMesh()
	vertices := []rendering.Vertex{
		{Pos: [3]float32{-1, -1, +0.0}},
		{Pos: [3]float32{+1, -1, +0.0}},
		{Pos: [3]float32{+0.0, +1, +0.0}},
	}
	mesh.AddVertices(vertices)

	meshRenderer := components.NewMeshRenderer(mesh)
	triangleObject := NewGameObject()
	triangleObject.AddComponent(meshRenderer)
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

func (g *Game) Update() {

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
