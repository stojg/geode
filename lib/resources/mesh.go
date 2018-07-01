package resources

import (
	"github.com/go-gl/gl/v4.1-core/gl"
)

// @todo check http://ogldev.atspace.co.uk/www/tutorial33/tutorial33.html for proper instanced rendering
// @todo Would be nice to have Model that has one or many meshes and textures
func NewMesh() *Mesh {
	m := &Mesh{}

	return m
}

type Mesh struct {
	vbo uint32
	vao uint32
	ebo uint32
	num int32
}

func (m *Mesh) SetVertices(vertices []Vertex, indices []uint32) {

	// Create buffers/arrays
	gl.GenBuffers(1, &m.vbo)
	gl.GenVertexArrays(1, &m.vao)
	gl.GenBuffers(1, &m.ebo)
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
	offset += numVertexPositions

	// normals
	gl.VertexAttribPointer(1, numVertexNormals, gl.FLOAT, false, int32(sizeOfVertex), gl.PtrOffset(offset*sizeOfFloat32))
	offset += numVertexNormals

	// texture coordinates
	gl.VertexAttribPointer(2, numVertexTexCoords, gl.FLOAT, false, int32(sizeOfVertex), gl.PtrOffset(offset*sizeOfFloat32))
	offset += numVertexTexCoords

	// tangents
	gl.VertexAttribPointer(3, numVertexTangents, gl.FLOAT, false, int32(sizeOfVertex), gl.PtrOffset(offset*sizeOfFloat32))

	// reset the current bound vertex array so that no one else mistakenly changes the VAO
	gl.BindVertexArray(0)
}

func (m *Mesh) Bind() {
	gl.BindVertexArray(m.vao)
	gl.EnableVertexAttribArray(0)
	gl.EnableVertexAttribArray(1)
	gl.EnableVertexAttribArray(2)
	gl.EnableVertexAttribArray(3)
}

func (m *Mesh) Draw() {
	gl.DrawElements(gl.TRIANGLES, m.num, gl.UNSIGNED_INT, gl.PtrOffset(0))
}

func (m *Mesh) Unbind() {
	gl.DisableVertexAttribArray(0)
	gl.DisableVertexAttribArray(1)
	gl.DisableVertexAttribArray(2)
	gl.DisableVertexAttribArray(3)
	gl.BindVertexArray(0)
}

func (m *Mesh) CleanUp() {
	//gl.DeleteV
}
