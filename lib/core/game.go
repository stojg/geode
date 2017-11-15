package core

import (
	"fmt"

	"github.com/go-gl/glfw/v3.2/glfw"
)

func NewGame() *Game {
	return &Game{}
}

type Game struct {}

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

}



