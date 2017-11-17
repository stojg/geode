package core

import (
	"fmt"

	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/stojg/graphics/lib/components"
	"github.com/stojg/graphics/lib/input"
	"github.com/stojg/graphics/lib/rendering"
)

type Drawable interface {
	Draw()
}

func NewGame() *Game {
	mesh := rendering.NewMesh()
	vertices := []rendering.Vertex{
		{Pos: [3]float32{-0.5, -0.5, +0.0}},
		{Pos: [3]float32{+0.5, -0.5, +0.0}},
		{Pos: [3]float32{+0.0, +0.5, +0.0}},
	}
	mesh.AddVertices(vertices)

	meshRenderer := components.NewMeshRenderer(mesh)

	g := &Game{}

	object := NewGameObject()
	object.AddComponent(meshRenderer)
	g.RootObject().AddChild(object)

	return g

}

type Game struct {
	root   *GameObject
	shader *rendering.Shader
}

func (g *Game) Input() {

	if input.KeyDown(glfw.KeyUp) {
		fmt.Println("up was just pressed")
	}

	if input.KeyUp(glfw.KeyUp) {
		fmt.Println("up was just relased")
	}

	if input.ButtonDown(glfw.MouseButton1) {
		fmt.Println("mouse 1 click")
		fmt.Println(input.CursorPosition())
	}

	if input.ButtonUp(glfw.MouseButton1) {
		fmt.Println("mouse 1 release")
		fmt.Println(input.CursorPosition())
	}
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
