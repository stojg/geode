package core

import (
	"fmt"
	"github.com/go-gl/gl/v4.1-core/gl"
)

func ClearScreen() {
	// @todo Stencil buffer
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
}

func InitGraphics() {

	if err := gl.Init(); err != nil {
		panic(err)
	}

	gl.ClearColor(0.05, 0.06, 0.07, 0)

	gl.FrontFace(gl.CCW)
	gl.CullFace(gl.BACK)
	gl.Enable(gl.CULL_FACE)
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)

	// @todo depth clamp for later

	gl.Enable(gl.FRAMEBUFFER_SRGB)
}

// CheckForError will check for OpenGL errors and panic if an error has been raised
func CheckForError(name string) {
	err := gl.GetError()
	switch err {
	case 0:
		return
	case gl.INVALID_OPERATION:
		fmt.Printf("GL Error: INVALID_OPERATION 0x0%x\n", err)
	case gl.INVALID_ENUM:
		fmt.Printf("GL Error: INVALID_ENUM 0x0%x\n", err)
	case gl.INVALID_FRAMEBUFFER_OPERATION:
		fmt.Printf("GL Error: INVALID_FRAMEBUFFER_OPERATION 0x0%x\n", err)
	default:
		fmt.Printf("GL Error: 0x0%x\n", err)
	}
}
