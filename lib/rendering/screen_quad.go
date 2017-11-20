package rendering

import "github.com/go-gl/gl/v4.1-core/gl"

func NewScreenQuad() *ScreenQuad {

	s := &ScreenQuad{}
	// vertex attributes for a quad that fills the entire screen in Normalized Device Coordinates.
	quadVertices := []float32{
		// positions // texCoords
		-1.0, 1.0, 0.0, 1.0,
		-1.0, -1.0, 0.0, 0.0,
		1.0, -1.0, 1.0, 0.0,

		-1.0, 1.0, 0.0, 1.0,
		1.0, -1.0, 1.0, 0.0,
		1.0, 1.0, 1.0, 1.0,
	}

	var vbo uint32
	gl.GenVertexArrays(1, &s.vao)
	gl.GenBuffers(1, &vbo)
	gl.BindVertexArray(s.vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(quadVertices)*sizeOfFloat32, gl.Ptr(quadVertices), gl.STATIC_DRAW)
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 2, gl.FLOAT, false, 4*4, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(1)
	gl.VertexAttribPointer(1, 2, gl.FLOAT, false, 4*4, gl.PtrOffset(2*sizeOfFloat32))

	return s
}

type ScreenQuad struct {
	vao uint32
}

func (s *ScreenQuad) Draw() {
	gl.BindVertexArray(s.vao)
	gl.DrawArrays(gl.TRIANGLES, 0, 6)
}
