package input

import (
	"github.com/go-gl/glfw/v3.2/glfw"
)

var (
	currentKeys map[glfw.Key]bool
	downKeys    map[glfw.Key]bool
	upKeys      map[glfw.Key]bool

	currentButtons map[glfw.MouseButton]bool
	downButtons    map[glfw.MouseButton]bool
	upButtons      map[glfw.MouseButton]bool

	// callback variables
	keys         [glfw.KeyLast]bool
	mouseButtons [glfw.MouseButtonLast]bool
	cursor       [2]float32
)

var inst *glfw.Window

func SetWindow(window *glfw.Window) {

	inst = window
	currentKeys = make(map[glfw.Key]bool)
	downKeys = make(map[glfw.Key]bool)
	upKeys = make(map[glfw.Key]bool)

	currentButtons = make(map[glfw.MouseButton]bool)
	downButtons = make(map[glfw.MouseButton]bool)
	upButtons = make(map[glfw.MouseButton]bool)

	window.SetKeyCallback(func(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		if action == glfw.Press {
			keys[key] = true
		} else if action == glfw.Release {
			keys[key] = false
		}
	})

	window.SetInputMode(glfw.CursorMode, glfw.CursorNormal)
	cX, xY := window.GetCursorPos()
	cursor[0], cursor[1] = float32(cX), float32(xY)
	window.SetCursorPosCallback(func(w *glfw.Window, xpos float64, ypos float64) {
		cursor[0], cursor[1] = float32(xpos), float32(ypos)
	})

	window.SetMouseButtonCallback(func(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mod glfw.ModifierKey) {
		if action == glfw.Press {
			mouseButtons[button] = true
		} else if action == glfw.Release {
			mouseButtons[button] = false
		}
	})
}

func Update() {
	for k := range currentKeys {
		upKeys[k] = !keys[k] && currentKeys[k]
		downKeys[k] = keys[k] && !currentKeys[k]
	}
	for k, v := range keys {
		currentKeys[glfw.Key(k)] = v
	}

	for k := range currentButtons {
		upButtons[k] = !mouseButtons[k] && currentButtons[k]
		downButtons[k] = mouseButtons[k] && !currentButtons[k]
	}
	for k, v := range mouseButtons {
		currentButtons[glfw.MouseButton(k)] = v
	}
}

func ShowCursor() {
	inst.SetInputMode(glfw.CursorMode, glfw.CursorNormal)
}

func HideCursor() {
	inst.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)
}

func SetCursorPosition(x, y float32) {
	inst.SetCursorPos(float64(x), float64(y))
	cursor[0], cursor[1] = x, y
}

func Key(keyCode glfw.Key) bool {
	return keys[keyCode]
}

func KeyDown(keyCode glfw.Key) bool {
	return downKeys[keyCode]
}

func KeyUp(keyCode glfw.Key) bool {
	return upKeys[keyCode]
}

func Button(button glfw.MouseButton) bool {
	return mouseButtons[button]
}

func ButtonDown(button glfw.MouseButton) bool {
	return downButtons[button]
}

func ButtonUp(button glfw.MouseButton) bool {
	return upButtons[button]
}

func CursorPosition() [2]float32 {
	return cursor
}
