package core

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/stojg/graphics/lib/components"

	"fmt"
	"runtime"
)

func NewWindow(width, height int, title string) (*Window, error) {
	w := &Window{
		width:  width,
		height: height,
		title:  title,
	}

	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()

	if err := glfw.Init(); err != nil {
		return w, fmt.Errorf("failed to initialize glfw: %s", err)
	}

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	glfw.WindowHint(glfw.SRGBCapable, glfw.True)
	//glfw.WindowHint(glfw.Samples, 0)

	// request a window
	window, err := glfw.CreateWindow(w.width, w.height, w.title, nil, nil)
	if err != nil {
		return w, err
	}
	w.inst = window

	w.inst.MakeContextCurrent()
	// disable or enable vertical refresh (vsync)
	glfw.SwapInterval(1)

	// the actual size of the window might be different due to screen, for example retina screens
	w.viewPortWidth, w.viewPortHeight = window.GetFramebufferSize()
	components.Width = w.viewPortWidth
	components.Height = w.viewPortHeight

	if err := gl.Init(); err != nil {
		return w, err
	}
	gl.Enable(gl.MULTISAMPLE)
	return w, nil
}

type Window struct {
	inst           *glfw.Window
	width          int
	height         int
	title          string
	viewPortHeight int
	viewPortWidth  int
}

func (w *Window) Instance() *glfw.Window {
	return w.inst
}

func (w *Window) Close() {
	glfw.Terminate()
}

func (w *Window) ShouldClose() bool {
	return w.inst.ShouldClose()
}

func (w *Window) Render() {
	w.inst.SwapBuffers()
	glfw.PollEvents()
}
