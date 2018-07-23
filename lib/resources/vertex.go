package resources

import "unsafe"

// number of elements in the Vertex.Pos
const numVertexPositions = 3
const numVertexNormals = 3
const numVertexTexCoords = 2
const numVertexTangents = 3

/* #nosec */
const sizeOfVertex = unsafe.Sizeof(Vertex{})

type Vertex struct {
	Pos       [3]float32
	Normal    [3]float32
	TexCoords [2]float32
	Tangent   [3]float32
}
