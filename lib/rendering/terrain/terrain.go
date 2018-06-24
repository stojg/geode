package terrain

import (
	"github.com/stojg/graphics/lib/rendering"
)

const Size float32 = 800
const VertexCount int = 128

func New(gridX, gridZ float32) *Terrain {
	data := generateTerrain()
	mesh := rendering.NewMesh()
	mesh.SetVertices(rendering.ConvertToVertices(data))

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

func generateTerrain() []float32 {
	// https://www.youtube.com/watch?v=yNYwZMmgTJk&list=PLRIWtICgwaX0u7Rf9zkZhLoLuZVfUksDP&index=14
	const count = VertexCount * VertexCount
	var vertices [count * 3]float32
	var normals [count * 3]float32
	var textureCoords [count * 2]float32

	vertexPointer := 0
	for i := 0; i < VertexCount; i++ {
		for j := 0; j < VertexCount; j++ {
			vertices[vertexPointer*3] = float32(j) / (float32(VertexCount) - 1) * Size
			vertices[vertexPointer*3+1] = 0
			vertices[vertexPointer*3+2] = float32(i) / (float32(VertexCount) - 1) * Size
			normals[vertexPointer*3] = 0
			normals[vertexPointer*3+1] = 1
			normals[vertexPointer*3+2] = 0
			textureCoords[vertexPointer*2] = float32(j) / (float32(VertexCount) - 1)
			textureCoords[vertexPointer*2+1] = float32(i) / (float32(VertexCount) - 1)
			vertexPointer++
		}
	}

	var indices [(VertexCount - 1) * (VertexCount - 1) * 6]int
	pointer := 0
	for gz := 0; gz < VertexCount-1; gz++ {
		for gx := 0; gx < VertexCount-1; gx++ {
			topLeft := (gz * VertexCount) + gx
			topRight := topLeft + 1
			bottomLeft := ((gz + 1) * VertexCount) + gx
			bottomRight := bottomLeft + 1
			indices[pointer] = topLeft
			pointer++
			indices[pointer] = bottomLeft
			pointer++
			indices[pointer] = topRight
			pointer++
			indices[pointer] = topRight
			pointer++
			indices[pointer] = bottomLeft
			pointer++
			indices[pointer] = bottomRight
			pointer++
		}
	}

	var result []float32
	for _, i := range indices {
		result = append(result, vertices[i*3])
		result = append(result, vertices[i*3+1])
		result = append(result, vertices[i*3+2])
		result = append(result, normals[i*3])
		result = append(result, normals[i*3+1])
		result = append(result, normals[i*3+2])
		result = append(result, textureCoords[i*2])
		result = append(result, textureCoords[i*2+1])
	}
	return result
}
