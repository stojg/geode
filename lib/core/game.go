package core

import (
	"fmt"
	"github.com/go-gl/glfw/v3.2/glfw"
)

func NewGame(s *Shader) *Game {
	g := &Game{
		mesh:   NewMesh(),
		shader: s,
	}

	vertices := []float32{
		-0.5, -0.5, 0.0,
		0.5, -0.5, 0.0,
		0.0, 0.5, 0.0,
	}
	g.mesh.AddVertices(vertices)
	return g
}

type Game struct {
	mesh   *Mesh
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
