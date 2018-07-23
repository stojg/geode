package primitives

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/stojg/graphics/lib/utilities"
)

var quadVao uint32 = 1<<32 - 1

func DrawQuad() {
	if quadVao == 1<<32-1 {
		setupQuad()
	}
	gl.BindVertexArray(quadVao)
	gl.EnableVertexAttribArray(0)
	gl.EnableVertexAttribArray(1)
	gl.DrawArrays(gl.TRIANGLE_STRIP, 0, 4)
	gl.BindVertexArray(0)

}

func setupQuad() {
	// vertex attributes for a quad that fills the entire screen in Normalized Device Coordinates.
	quadVertices := []float32{
		// positions -  textureCoords
		-1.0, 1.0, 0.0, 0.0, 1.0,
		-1.0, -1.0, 0.0, 0.0, 0.0,
		1.0, 1.0, 0.0, 1.0, 1.0,
		1.0, -1.0, 0.0, 1.0, 0.0,
	}
	var vbo uint32
	gl.GenVertexArrays(1, &quadVao)
	gl.BindVertexArray(quadVao)

	vbo = utilities.CreateEmptyVBO(len(quadVertices), gl.STATIC_DRAW)
	utilities.UpdateVBO(vbo, len(quadVertices), quadVertices, gl.STATIC_DRAW)
	instanceDataLength := 5 // 3 pos, 2 TexCoords
	utilities.AddAttribute(quadVao, vbo, 0, 3, instanceDataLength, 0)
	utilities.AddAttribute(quadVao, vbo, 1, 2, instanceDataLength, 3)

}
