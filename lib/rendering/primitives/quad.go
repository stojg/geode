package primitives

import "github.com/go-gl/gl/v4.1-core/gl"

var quadVao uint32 = 1<<32 - 1

func DrawQuad() {
	if quadVao == 1<<32-1 {
		setupQuad()
	}
	gl.BindVertexArray(quadVao)
	gl.DrawArrays(gl.TRIANGLE_STRIP, 0, 4)

}

func setupQuad() {
	// vertex attributes for a quad that fills the entire screen in Normalized Device Coordinates.
	quadVertices := []float32{
		// positions -  texture Coords
		-1.0, 1.0, 0.0, 0.0, 1.0,
		-1.0, -1.0, 0.0, 0.0, 0.0,
		1.0, 1.0, 0.0, 1.0, 1.0,
		1.0, -1.0, 0.0, 1.0, 0.0,
	}
	var vbo uint32
	gl.GenVertexArrays(1, &quadVao)
	gl.GenBuffers(1, &vbo)
	gl.BindVertexArray(quadVao)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(quadVertices)*sizeOfFloat32, gl.Ptr(quadVertices), gl.STATIC_DRAW)
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 5*4, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(1)
	gl.VertexAttribPointer(1, 2, gl.FLOAT, false, 5*4, gl.PtrOffset(3*sizeOfFloat32))
}
