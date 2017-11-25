package core

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"

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
	glfw.WindowHint(glfw.Samples, 4)
	glfw.WindowHint(glfw.SRGBCapable, glfw.True)

	// request a window
	window, err := glfw.CreateWindow(w.width, w.height, w.title, nil, nil)
	if err != nil {
		return w, err
	}
	w.inst = window

	w.inst.MakeContextCurrent()
	// disable or enable vertical refresh (vsync)
	glfw.SwapInterval(0)

	// this is the actual numVertices we got
	fbw, fbh := window.GetFramebufferSize()
	w.viewPortWidth = int32(fbw)
	w.viewPortHeight = int32(fbh)

	return w, gl.Init()
}

type Window struct {
	inst           *glfw.Window
	width          int
	height         int
	title          string
	viewPortHeight int32
	viewPortWidth  int32
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
