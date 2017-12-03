package rendering

import (
	"fmt"

	"github.com/go-gl/gl/v4.1-core/gl"
)

// checkForError will check for OpenGL errors and panic if an error has been raised
func checkForError(name string) bool {
	err := gl.GetError()
	switch err {
	case 0:
		return false
	case gl.INVALID_OPERATION:
		fmt.Printf("[%s] GL Error: INVALID_OPERATION 0x0%x\n", name, err)
	case gl.INVALID_ENUM:
		fmt.Printf("[%s] GL Error: INVALID_ENUM 0x0%x\n", name, err)
	case gl.INVALID_FRAMEBUFFER_OPERATION:
		fmt.Printf("[%s] GL Error: INVALID_FRAMEBUFFER_OPERATION 0x0%x\n", name, err)
	default:
		fmt.Printf("[%s] GL Error: 0x0%x\n", name, err)
	}
	return true
}
