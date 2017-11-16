package core

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

func (m *Mesh) AddVertices(vertices []float32) {
	const sizeFloat = 4
	const VertexSize = 3

	m.numVertices = int32(len(vertices) / VertexSize)

	gl.BindVertexArray(m.vao)

	// load data into vertex buffers
	gl.BindBuffer(gl.ARRAY_BUFFER, m.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*sizeFloat, gl.Ptr(vertices), gl.STATIC_DRAW)

	// vertex position
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, VertexSize, gl.FLOAT, false, VertexSize*sizeFloat, gl.PtrOffset(0))

	// @todo add vertex attribute pointers for normals, texture coordinates and tangents

	// reset the current bound vertex array so that no one else mistakenly changes the VAO
	gl.BindVertexArray(0)
}

func (m *Mesh) Draw() {
	gl.BindVertexArray(m.vao)
	gl.DrawArrays(gl.TRIANGLES, 0, m.numVertices)
	CheckForError("Mesh.draw [end]")
}
