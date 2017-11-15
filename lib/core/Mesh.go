package core

import (

	"github.com/go-gl/gl/v4.1-core/gl"
)

func NewMesh() *Mesh {
	m := &Mesh{}
	gl.GenBuffers(1, &m.vbo)
	return m
}

type Mesh struct {
	vbo uint32
	size int32
}

func (m *Mesh) AddVertices(vertices []Vertex) {
	m.size = int32(len(vertices))
	gl.BindBuffer(gl.ARRAY_BUFFER, m.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, int(m.size * VertexSize), gl.Ptr(vertices), gl.STATIC_DRAW)
}

func (m *Mesh) Draw() {
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 4 * VertexSize, nil)
	gl.DrawArrays(gl.TRIANGLES, 0, m.size)
}
