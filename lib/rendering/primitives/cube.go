package primitives

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/stojg/graphics/lib/utilities"
)

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
	vertices := []float32{1, -1, -1, 1, -1, 1, -1, -1, 1, -1, -1, -1, 1, 1, -1, -1, 1, -1, -1, 1, 1, 1, 1, 1}
	indices := []uint32{0, 1, 2, 0, 2, 3, 4, 5, 6, 4, 6, 7, 0, 4, 7, 0, 7, 1, 1, 7, 6, 1, 6, 2, 2, 6, 5, 2, 5, 3, 4, 0, 3, 4, 3, 5}

	gl.GenVertexArrays(1, &cubeVao)
	vbo := utilities.CreateFloatVBO(cubeVao, len(vertices), vertices, gl.STATIC_DRAW)
	utilities.CreateIntEBO(cubeVao, len(indices), indices, gl.STATIC_DRAW)
	utilities.AddAttribute(cubeVao, vbo, 0, 3, 3, 0)
}
