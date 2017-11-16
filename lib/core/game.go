package core

import (
	"fmt"
	"github.com/go-gl/glfw/v3.2/glfw"
)

type Drawable interface{
	Draw()
}

func NewGame(s *Shader) *Game {
	mesh := NewMesh()
	vertices := []float32{
		-0.5, -0.5, 0.0,
		0.5, -0.5, 0.0,
		0.0, 0.5, 0.0,
	}
	mesh.AddVertices(vertices)

	return &Game{
		mesh:   mesh,
		shader: s,
	}
}

type Game struct {
	mesh   Drawable
	shader *Shader
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
	g.shader.Bind()
	g.mesh.Draw()
}
