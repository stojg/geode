package core

import (
	"fmt"

	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

func NewGame() *Game {
	g := &Game{
		mesh: NewMesh(),
	}

	vertices := []Vertex{
		NewVertex(mgl32.Vec3{-1, -1, 0}),
		NewVertex(mgl32.Vec3{0, 1, 0}),
		NewVertex(mgl32.Vec3{-1, 1, 0}),
	}
	g.mesh.AddVertices(vertices)
	return g
}

type Game struct {
	mesh *Mesh
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
		fmt.Println(i.MousePosition())
	}

	if i.ButtonUp(glfw.MouseButton1) {
		fmt.Println("mouse 1 release")
		fmt.Println(i.MousePosition())
	}
}

func (g *Game) Update() {

}

func (g *Game) Render() {
	g.mesh.Draw()
}
