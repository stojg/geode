package rendering

import (
	"fmt"

	"github.com/go-gl/gl/v4.1-core/gl"
)

// checkForError will check for OpenGL errors and panic if an error has been raised
func checkForError(name string) {
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
