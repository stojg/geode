package primitives

import (
	"unsafe"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

/* #nosec */
const SizeOfFloat32 = int(unsafe.Sizeof(float32(1)))
const SizeOfMat4 = int(unsafe.Sizeof(mgl32.Ident4()))

var cubeVao uint32 = 1<<32 - 1

func DrawCube() {
	if cubeVao == 1<<32-1 {
		setupCube()
	}
	gl.BindVertexArray(cubeVao)
	gl.EnableVertexAttribArray(0)
	gl.DrawElements(gl.TRIANGLES, 36, gl.UNSIGNED_INT, gl.PtrOffset(0))
	gl.BindVertexArray(0)
}

func setupCube() {
	inds := []uint32{0, 1, 2, 0, 2, 3, 4, 5, 6, 4, 6, 7, 0, 4, 7, 0, 7, 1, 1, 7, 6, 1, 6, 2, 2, 6, 5, 2, 5, 3, 4, 0, 3, 4, 3, 5}
	verts := []float32{1, -1, -1, 1, -1, 1, -1, -1, 1, -1, -1, -1, 1, 1, -1, -1, 1, -1, -1, 1, 1, 1, 1, 1}

	var vbo uint32
	var ebo uint32
	gl.GenVertexArrays(1, &cubeVao)

	gl.GenBuffers(1, &vbo)
	gl.GenBuffers(1, &ebo)

	gl.BindVertexArray(cubeVao)

	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(verts)*SizeOfFloat32, gl.Ptr(verts), gl.STATIC_DRAW)

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(inds)*4, gl.Ptr(inds), gl.STATIC_DRAW)

	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, int32(3*SizeOfFloat32), gl.PtrOffset(0))

	gl.BindVertexArray(0)
}
