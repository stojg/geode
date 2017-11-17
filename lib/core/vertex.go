package core

import "unsafe"

// number of elements in the Vertex.Pos
const numVertexPositions = 3

const sizeOfVertex = unsafe.Sizeof(Vertex{})

const sizeOfFloat32 = int(unsafe.Sizeof(float32(1)))

type Vertex struct {
	Pos [3]float32
}
