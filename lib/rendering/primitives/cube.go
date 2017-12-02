package primitives

import (
	"unsafe"

	"github.com/go-gl/gl/v4.1-core/gl"
)

/* #nosec */
const sizeOfFloat32 = int(unsafe.Sizeof(float32(1)))

var cubeVao uint32 = 1<<32 - 1

func DrawCube() {
	if cubeVao == 1<<32-1 {
		setupCube()
	}
	gl.BindVertexArray(cubeVao)
	gl.DrawArrays(gl.TRIANGLES, 0, 36)

}

func setupCube() {
	quadVertices := []float32{
		1, -1, -1,
		1, -1, 1,
		-1, -1, 1,
		1, -1, -1,
		-1, -1, 1,
		-1, -1, -1,
		1, 1, -1,
		-1, 1, -1,
		-1, 1, 1,
		1, 1, -1,
		-1, 1, 1,
		1, 1, 1,
		1, -1, -1,
		1, 1, -1,
		1, 1, 1,
		1, -1, -1,
		1, 1, 1,
		1, -1, 1,
		1, -1, 1,
		1, 1, 1,
		-1, 1, 1,
		1, -1, 1,
		-1, 1, 1,
		-1, -1, 1,
		-1, -1, 1,
		-1, 1, 1,
		-1, 1, -1,
		-1, -1, 1,
		-1, 1, -1,
		-1, -1, -1,
		1, 1, -1,
		1, -1, -1,
		-1, -1, -1,
		1, 1, -1,
		-1, -1, -1,
		-1, 1, -1,
	}

	var vbo uint32
	gl.GenVertexArrays(1, &cubeVao)
	gl.GenBuffers(1, &vbo)
	gl.BindVertexArray(cubeVao)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(quadVertices)*sizeOfFloat32, gl.Ptr(quadVertices), gl.STATIC_DRAW)
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 4*3, gl.PtrOffset(0))
}
