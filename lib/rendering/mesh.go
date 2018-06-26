package rendering

import (
	"github.com/go-gl/gl/v4.1-core/gl"
)

func NewMesh() *Mesh {
	m := &Mesh{}

	// Create buffers/arrays
	gl.GenBuffers(1, &m.vbo)
	gl.GenVertexArrays(1, &m.vao)
	gl.GenBuffers(1, &m.ebo)

	return m
}

type Mesh struct {
	vbo uint32
	vao uint32
	ebo uint32
	num int32
}

func (m *Mesh) SetVertices(vertices []Vertex, indices []uint32) {

	m.num = int32(len(indices))

	gl.BindVertexArray(m.vao)

	// load data into vertex buffers
	gl.BindBuffer(gl.ARRAY_BUFFER, m.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*int(sizeOfVertex), gl.Ptr(vertices), gl.STATIC_DRAW)

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, m.ebo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*int(sizeOfUint32), gl.Ptr(indices), gl.STATIC_DRAW)

	offset := 0
	// vertex position
	gl.VertexAttribPointer(0, numVertexPositions, gl.FLOAT, false, int32(sizeOfVertex), gl.PtrOffset(offset*sizeOfFloat32))
	gl.EnableVertexAttribArray(0)
	offset += numVertexPositions

	// normals
	gl.VertexAttribPointer(1, numVertexNormals, gl.FLOAT, false, int32(sizeOfVertex), gl.PtrOffset(offset*sizeOfFloat32))
	gl.EnableVertexAttribArray(1)
	offset += numVertexNormals

	// texture coordinates
	gl.VertexAttribPointer(2, numVertexTexCoords, gl.FLOAT, false, int32(sizeOfVertex), gl.PtrOffset(offset*sizeOfFloat32))
	gl.EnableVertexAttribArray(2)
	offset += numVertexTexCoords

	// tangents
	gl.VertexAttribPointer(3, 3, gl.FLOAT, false, int32(sizeOfVertex), gl.PtrOffset(offset*sizeOfFloat32))
	gl.EnableVertexAttribArray(3)

	// reset the current bound vertex array so that no one else mistakenly changes the VAO
	gl.BindVertexArray(0)
}

func (m *Mesh) Prepare() {
	gl.BindVertexArray(m.vao)
}

func (m *Mesh) Draw() {
	gl.DrawElements(gl.TRIANGLES, m.num, gl.UNSIGNED_INT, gl.PtrOffset(0))
}
