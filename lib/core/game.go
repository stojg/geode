package core

import (
	"fmt"
	"github.com/go-gl/glfw/v3.2/glfw"
)

type Drawable interface {
	Draw()
}

type GameObject interface {
	Input(*Input)
	Update()
	Render(*Shader)
}

func NewGame(s *Shader) *Game {
	mesh := NewMesh()
	vertices := []Vertex{
		{Pos: [3]float32{-0.5, -0.5, +0.0}},
		{Pos: [3]float32{+0.5, -0.5, +0.0}},
		{Pos: [3]float32{+0.0, +0.5, +0.0}},
	}
	mesh.AddVertices(vertices)

	meshRenderer := NewMeshRenderer(mesh)

	return &Game{
		gameObjects: []GameObject{meshRenderer},
		shader:      s,
	}
}

type Game struct {
	gameObjects []GameObject
	shader      *Shader
}

func (g *Game) Input(i *Input) {

	if i.KeyDown(glfw.KeyUp) {
		fmt.Println("up was just pressed")
	}

	if i.KeyUp(glfw.KeyUp) {
		fmt.Println("up was just relased")
	}

	if i.ButtonDown(glfw.MouseButton1) {
		fmt.Println("mouse 1 click")
		fmt.Println(i.CursorPosition())
	}

	if i.ButtonUp(glfw.MouseButton1) {
		fmt.Println("mouse 1 release")
		fmt.Println(i.CursorPosition())
	}
}

func (g *Game) Update() {

}

func (g *Game) Render() {
	for _, object := range g.gameObjects {
		object.Render(g.shader)
	}
}
