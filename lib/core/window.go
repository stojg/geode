package core

import (
	//"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"

	"fmt"
	"runtime"
)

var keys [glfw.KeyLast]bool
var mouseButtons [glfw.MouseButtonLast]bool
var cursorPosition mgl32.Vec2

func NewWindow(width, height int, title string) *Window {
	w := &Window{
		width: width,
		height: height,
		title: title,
	}
	w.Open()
	return w
}

type Window struct {
	inst   *glfw.Window
	width  int
	height int
	title  string
	viewPortHeight int32
	viewPortWidth int32

	previousFrameSec float64
	frameCounter     int
}

func (w *Window) Open() error {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()

	if err := glfw.Init(); err != nil {
		return fmt.Errorf("failed to initialize glfw: %s", err)
	}

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	glfw.WindowHint(glfw.Samples, 0)

	// request a window
	window, err := glfw.CreateWindow(w.width, w.height, w.title, nil, nil)
	if err != nil {
		return err
	}
	w.inst = window

	w.inst.MakeContextCurrent()
	// disable or enable vertical refresh (vsync)
	glfw.SwapInterval(0)

	window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)
	window.SetKeyCallback(func(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		if action == glfw.Press {
			keys[key] = true
		} else if action == glfw.Release {
			keys[key] = false
		}
	})
	window.SetMouseButtonCallback(func(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mod glfw.ModifierKey){
		if action == glfw.Press {
			mouseButtons[button] = true
		} else if action == glfw.Release {
			mouseButtons[button] = false
		}
	})

	x, y := window.GetCursorPos()
	cursorPosition[0] = float32(x)
	cursorPosition[1] = float32(y)
	window.SetCursorPosCallback(func(w *glfw.Window, xpos float64, ypos float64){
		cursorPosition[0] = float32(xpos)
		cursorPosition[1] = float32(ypos)
	})

	// this is the actual size we got
	fbw, fbh := window.GetFramebufferSize()
	w.viewPortWidth = int32(fbw)
	w.viewPortHeight = int32(fbh)

	return nil
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
