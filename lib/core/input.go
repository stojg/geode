package core

import (
	"github.com/go-gl/glfw/v3.2/glfw"
)

func NewInput(window *glfw.Window) *Input {

	input := &Input{
		currentKeys: make(map[glfw.Key]bool),
		downKeys:    make(map[glfw.Key]bool),
		upKeys:      make(map[glfw.Key]bool),

		currentButtons: make(map[glfw.MouseButton]bool),
		downButtons:    make(map[glfw.MouseButton]bool),
		upButtons:      make(map[glfw.MouseButton]bool),
	}

	window.SetKeyCallback(func(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		if action == glfw.Press {
			input.keys[key] = true
		} else if action == glfw.Release {
			input.keys[key] = false
		}
	})

	window.SetInputMode(glfw.CursorMode, glfw.CursorNormal)
	input.cursor[0], input.cursor[1] = window.GetCursorPos()
	window.SetCursorPosCallback(func(w *glfw.Window, xpos float64, ypos float64) {
		input.cursor[0], input.cursor[1] = xpos, ypos
	})

	window.SetMouseButtonCallback(func(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mod glfw.ModifierKey) {
		if action == glfw.Press {
			input.mouseButtons[button] = true
		} else if action == glfw.Release {
			input.mouseButtons[button] = false
		}
	})

	return input
}

type Input struct {
	currentKeys map[glfw.Key]bool
	downKeys    map[glfw.Key]bool
	upKeys      map[glfw.Key]bool

	currentButtons map[glfw.MouseButton]bool
	downButtons    map[glfw.MouseButton]bool
	upButtons      map[glfw.MouseButton]bool

	// callback variables
	keys         [glfw.KeyLast]bool
	mouseButtons [glfw.MouseButtonLast]bool
	cursor       [2]float64
}

func (input *Input) Update() {
	for k := range input.currentKeys {
		input.upKeys[k] = !input.keys[k] && input.currentKeys[k]
		input.downKeys[k] = input.keys[k] && !input.currentKeys[k]
	}
	for k, v := range input.keys {
		input.currentKeys[glfw.Key(k)] = v
	}

	for k := range input.currentButtons {
		input.upButtons[k] = !input.mouseButtons[k] && input.currentButtons[k]
		input.downButtons[k] = input.mouseButtons[k] && !input.currentButtons[k]
	}
	for k, v := range input.mouseButtons {
		input.currentButtons[glfw.MouseButton(k)] = v
	}
}

func (input *Input) Key(keyCode glfw.Key) bool {
	return input.keys[keyCode]
}

func (input *Input) KeyDown(keyCode glfw.Key) bool {
	return input.downKeys[keyCode]
}

func (input *Input) KeyUp(keyCode glfw.Key) bool {
	return input.upKeys[keyCode]
}

func (input *Input) Button(button glfw.MouseButton) bool {
	return input.mouseButtons[button]
}

func (input *Input) ButtonDown(button glfw.MouseButton) bool {
	return input.downButtons[button]
}

func (input *Input) ButtonUp(button glfw.MouseButton) bool {
	return input.upButtons[button]
}

func (input *Input) CursorPosition() [2]float64 {
	return input.cursor
}
