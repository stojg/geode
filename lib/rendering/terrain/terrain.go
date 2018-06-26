package terrain

import (
	"github.com/stojg/graphics/lib/rendering"
)

const Size float32 = 800
const VertexCount int = 128

func New(gridX, gridZ float32) *Terrain {
	v, i := generateTerrain()
	mesh := rendering.NewMesh()
	mesh.SetVertices(rendering.ConvertToVertices(v, i), i)

	return &Terrain{
		x:    gridX * Size,
		z:    gridZ * Size,
		mesh: mesh,
	}
}

type Terrain struct {
	x, z float32

	mesh    *rendering.Mesh
	texture interface{}
}

func (t *Terrain) Z() float32 {
	return t.z
}

func (t *Terrain) X() float32 {
	return t.x
}

func (t *Terrain) Mesh() *rendering.Mesh {
	return t.mesh
}

func generateTerrain() ([]float32, []uint32) {
	// https://www.youtube.com/watch?v=yNYwZMmgTJk&list=PLRIWtICgwaX0u7Rf9zkZhLoLuZVfUksDP&index=14
	var vertices []float32
	for i := 0; i < VertexCount; i++ {
		for j := 0; j < VertexCount; j++ {
			vertices = append(vertices, float32(j)/(float32(VertexCount)-1)*Size)
			vertices = append(vertices, 0)
			vertices = append(vertices, float32(i)/(float32(VertexCount)-1)*Size)
			vertices = append(vertices, 0)
			vertices = append(vertices, 1)
			vertices = append(vertices, 0)
			vertices = append(vertices, float32(j)/(float32(VertexCount)-1))
			vertices = append(vertices, float32(i)/(float32(VertexCount)-1))
		}
	}

	var indices []uint32
	for gz := 0; gz < VertexCount-1; gz++ {
		for gx := 0; gx < VertexCount-1; gx++ {
			topLeft := uint32((gz * VertexCount) + gx)
			topRight := uint32(topLeft + 1)
			bottomLeft := uint32(((gz + 1) * VertexCount) + gx)
			bottomRight := uint32(bottomLeft + 1)

			indices = append(indices, topLeft)
			indices = append(indices, bottomLeft)
			indices = append(indices, topRight)
			indices = append(indices, topRight)
			indices = append(indices, bottomLeft)
			indices = append(indices, bottomRight)
		}
	}

	return vertices, indices
}
