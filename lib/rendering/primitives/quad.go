package primitives

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/stojg/geode/lib/buffers"
	"github.com/stojg/geode/lib/debug"
)

var quadVao uint32 = 1<<32 - 1

func DrawQuad() {
	if quadVao == 1<<32-1 {
		setupQuad()
	}
	gl.BindVertexArray(quadVao)
	debug.AddVertexBind()
	gl.EnableVertexAttribArray(0)
	gl.EnableVertexAttribArray(1)
	gl.DrawArrays(gl.TRIANGLE_STRIP, 0, 4)
	gl.BindVertexArray(0)
	debug.AddVertexBind()

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

	vbo = buffers.CreateFloatVBO(quadVao, len(quadVertices), quadVertices, gl.STATIC_DRAW)
	instanceDataLength := int32(5) // 3 pos, 2 TexCoords
	buffers.AddAttribute(quadVao, vbo, 0, 3, instanceDataLength, 0)
	buffers.AddAttribute(quadVao, vbo, 1, 2, instanceDataLength, 3)

}
