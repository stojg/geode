package rendering

import (
	"github.com/go-gl/gl/v4.1-core/gl"
)

func NewMesh() *Mesh {
	m := &Mesh{}

	// Create buffers/arrays
	gl.GenBuffers(1, &m.vbo)
	gl.GenVertexArrays(1, &m.vao)

	return m
}

type Mesh struct {
	vbo         uint32
	vao         uint32
	numVertices int32
}

func (m *Mesh) AddVertices(vertices []Vertex) {

	m.numVertices = int32(len(vertices))

	gl.BindVertexArray(m.vao)

	// load data into vertex buffers
	gl.BindBuffer(gl.ARRAY_BUFFER, m.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*int(sizeOfVertex), gl.Ptr(vertices), gl.STATIC_DRAW)

	offset := 0
	// vertex position
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, numVertexPositions, gl.FLOAT, false, int32(sizeOfVertex), gl.PtrOffset(offset*sizeOfFloat32))
	offset += numVertexPositions

	// normals
	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, int32(sizeOfVertex), gl.PtrOffset(offset*sizeOfFloat32))
	gl.EnableVertexAttribArray(1)
	offset += numVertexNormals

	// texture coordinates
	gl.VertexAttribPointer(2, 2, gl.FLOAT, false, int32(sizeOfVertex), gl.PtrOffset(offset*sizeOfFloat32))
	gl.EnableVertexAttribArray(2)
	offset += numVertexTexCoords

	// tangents
	gl.VertexAttribPointer(3, 3, gl.FLOAT, false, int32(sizeOfVertex), gl.PtrOffset(offset*sizeOfFloat32))
	gl.EnableVertexAttribArray(3)

	// reset the current bound vertex array so that no one else mistakenly changes the VAO
	gl.BindVertexArray(0)
}

func (m *Mesh) Draw() {
	gl.BindVertexArray(m.vao)
	gl.DrawArrays(gl.TRIANGLES, 0, m.numVertices)
}
