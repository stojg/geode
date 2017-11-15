package core

import "github.com/go-gl/mathgl/mgl32"

const VertexSize = 3
func NewVertex(pos mgl32.Vec3) Vertex {
	return Vertex{
		Pos: pos,
	}
}

type Vertex struct {
	Pos mgl32.Vec3
}


