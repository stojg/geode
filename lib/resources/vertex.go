package resources

import "unsafe"

/* #nosec */
const sizeOfVertex = unsafe.Sizeof(Vertex{})

type Vertex struct {
	Pos       [3]float32
	Normal    [3]float32
	TexCoords [2]float32
	Tangent   [3]float32
}
