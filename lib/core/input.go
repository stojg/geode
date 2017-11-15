package core

import (
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

func NewInput() *Input {
	return &Input{
		currentKeys: make(map[glfw.Key]bool),
		downKeys: make(map[glfw.Key]bool),
		upKeys: make(map[glfw.Key]bool),

		currentButtons: make(map[glfw.MouseButton]bool),
		downButtons:    make(map[glfw.MouseButton]bool),
		upButtons:      make(map[glfw.MouseButton]bool),
	}
}

type Input struct {
	currentKeys map[glfw.Key]bool
	downKeys map[glfw.Key]bool
	upKeys map[glfw.Key]bool

	currentButtons map[glfw.MouseButton]bool
	downButtons    map[glfw.MouseButton]bool
	upButtons      map[glfw.MouseButton]bool
}

func (input *Input) Update() {
	for k := range input.currentKeys {
		input.upKeys[k] = !keys[k] && input.currentKeys[k]
		input.downKeys[k] = keys[k] && !input.currentKeys[k]
	}
	for k, v := range keys {
		input.currentKeys[glfw.Key(k)] = v
	}

	for k := range input.currentButtons {
		input.upButtons[k] = !mouseButtons[k] && input.currentButtons[k]
		input.downButtons[k] = mouseButtons[k] && !input.currentButtons[k]
	}
	for k, v := range mouseButtons {
		input.currentButtons[glfw.MouseButton(k)] = v
	}
}

func (input *Input) Key(keyCode glfw.Key) bool {
	return keys[keyCode]
}

func (input *Input) KeyDown(keyCode glfw.Key) bool {
	return input.downKeys[keyCode]
}

func (input *Input) KeyUp(keyCode glfw.Key) bool {
	return input.upKeys[keyCode]
}

func (input *Input) Button(button  glfw.MouseButton) bool {
	return mouseButtons[button]
}

func (input *Input) ButtonDown(button  glfw.MouseButton) bool {
	return input.downButtons[button]
}

func (input *Input) ButtonUp(button glfw.MouseButton) bool {
	return input.upButtons[button]
}

func (input *Input) MousePosition() mgl32.Vec2 {
	return cursorPosition
}
