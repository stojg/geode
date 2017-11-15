package core

import "github.com/go-gl/gl/v4.1-core/gl"

func ClearScreen() {
	// @todo Stencil buffer
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
}

func InitGraphics() {

	if err := gl.Init(); err != nil {
		panic(err)
	}

	gl.ClearColor(0,0,0,0)

	gl.FrontFace(gl.CW)
	gl.CullFace(gl.BACK)
	gl.Enable(gl.CULL_FACE)
	gl.Enable(gl.DEPTH_TEST)

	// @todo depth clamp for later

	gl.Enable(gl.FRAMEBUFFER_SRGB)
}
