package debug

import (
	"fmt"

	"github.com/go-gl/gl/v4.1-core/gl"
)

// CheckForError will check for OpenGL errors and return true if an error was reported
func CheckForError(debugName string) bool {
	err := gl.GetError()
	switch err {
	case 0:
		return false
	case gl.INVALID_OPERATION:
		fmt.Printf("[%s] GL Error: INVALID_OPERATION 0x0%x\n", debugName, err)
	case gl.INVALID_ENUM:
		fmt.Printf("[%s] GL Error: INVALID_ENUM 0x0%x\n", debugName, err)
	case gl.INVALID_FRAMEBUFFER_OPERATION:
		fmt.Printf("[%s] GL Error: INVALID_FRAMEBUFFER_OPERATION 0x0%x\n", debugName, err)
	default:
		fmt.Printf("[%s] GL Error: 0x0%x\n", debugName, err)
	}
	return true
}

func FramebufferComplete(debugName string) {
	if e := gl.CheckFramebufferStatus(gl.FRAMEBUFFER); e != gl.FRAMEBUFFER_COMPLETE {
		panic(fmt.Sprintf("%s Framebuffer creation failed, FBO isn't complete: 0x%x", debugName, e))
	}
}

var drawcalls uint64 = 0

func Drawcall() {
	drawcalls++
}

func GetDrawcalls() uint64 {
	t := drawcalls
	drawcalls = 0
	return t
}

var shaderSwitch uint64 = 0

func ShaderSwitch() {
	shaderSwitch++
}

func GetShaderSwitches() uint64 {
	t := shaderSwitch
	shaderSwitch = 0
	return t
}

var uniformSet uint64 = 0

func AddUniformSet() {
	uniformSet++
}

func GetUniformSet() uint64 {
	t := uniformSet
	uniformSet = 0
	return t
}

var vertexbind uint64 = 0

func AddVertexBind() {
	vertexbind++
}

func GetVertexBind() uint64 {
	t := vertexbind
	vertexbind = 0
	return t
}

var particles uint64

func SetParticles(num uint64) {
	particles = num
}

func GetParticles() uint64 {
	t := particles
	particles = 0
	return t
}
