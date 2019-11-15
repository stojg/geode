package terrain

import (
	"math"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/stojg/geode/lib/components"
	"github.com/stojg/geode/lib/resources"
)

const Size float32 = 512
const VertexCount = 128
const stride = 8

func New(gridX, gridZ float32) *Terrain {
	t := &Terrain{
		worldX:         gridX * Size,
		worldZ:         gridZ * Size,
		gridSizeSquare: Size / (float32(VertexCount) - 1),
	}

	v, i := t.generateTerrain(gridX, gridZ)
	mesh := resources.NewMesh()
	mesh.SetVertices(resources.ConvertToVertices(v, i), i)
	t.mesh = mesh
	return t
}

type Terrain struct {
	worldX, worldZ float32
	heights        [VertexCount][VertexCount]float32
	gridSizeSquare float32
	mesh           components.Drawable
}

func (t *Terrain) generateTerrain(gridX, gridZ float32) ([]float32, []uint32) {
	xOffset = float64(gridX*float32(VertexCount-1)) / 4
	zOffset = float64(gridZ*float32(VertexCount-1)) / 4

	hg := NewHeightGenerator(349)

	// https://www.youtube.com/watch?v=yNYwZMmgTJk&list=PLRIWtICgwaX0u7Rf9zkZhLoLuZVfUksDP&index=14
	vertices := make([]float32, VertexCount*VertexCount*stride)

	for z := 0; z < VertexCount; z++ {
		for x := 0; x < VertexCount; x++ {
			t.heights[x][z] = hg.Height(float64(x)+xOffset, float64(z)+zOffset)
			numSquares := float32(VertexCount) - 1
			v := []float32{
				float32(x) / numSquares * Size,
				t.heights[x][z],
				float32(z) / numSquares * Size,
				0, 0, 0,
				float32(x) / numSquares,
				float32(z) / numSquares,
			}
			copy(vertices[x*stride+z*stride*VertexCount:], v)
		}
	}

	vertices = t.setNormals(vertices)
	//vertices = t.normaliseNormals(vertices)

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

func (t *Terrain) Height(x, z float32) float32 {
	terrainX := x - t.worldX
	terrainZ := z - t.worldZ

	gridX := int(math.Floor(float64(terrainX / t.gridSizeSquare)))
	gridZ := int(math.Floor(float64(terrainZ / t.gridSizeSquare)))

	if gridX >= VertexCount-1 || gridX < 0 {
		return 0
	}
	if gridZ >= VertexCount-1 || gridZ < 0 {
		return 0
	}

	xCoord := float32(math.Mod(float64(terrainX), float64(t.gridSizeSquare))) / t.gridSizeSquare
	zCoord := float32(math.Mod(float64(terrainZ), float64(t.gridSizeSquare))) / t.gridSizeSquare

	var result float32

	if xCoord <= 1-zCoord {
		result = barryCentric([3]float32{0, t.heights[gridX][gridZ], 0}, [3]float32{1, t.heights[gridX+1][gridZ], 0}, [3]float32{0, t.heights[gridX][gridZ+1], 1}, [2]float32{xCoord, zCoord})
	} else {
		result = barryCentric([3]float32{1, t.heights[gridX+1][gridZ], 0}, [3]float32{1, t.heights[gridX+1][gridZ+1], 1}, [3]float32{0, t.heights[gridX][gridZ+1], 1}, [2]float32{xCoord, zCoord})
	}

	return result
}

func (t *Terrain) Z() float32 {
	return t.worldZ
}

func (t *Terrain) X() float32 {
	return t.worldX
}

func (t *Terrain) Mesh() components.Drawable {
	return t.mesh
}

// @todo add to utilities
func barryCentric(p1, p2, p3 [3]float32, pos [2]float32) float32 {
	det := (p2[2]-p3[2])*(p1[0]-p3[0]) + (p3[0]-p2[0])*(p1[2]-p3[2])
	l1 := ((p2[2]-p3[2])*(pos[0]-p3[0]) + (p3[0]-p2[0])*(pos[1]-p3[2])) / det
	l2 := ((p3[2]-p1[2])*(pos[0]-p3[0]) + (p1[0]-p3[0])*(pos[1]-p3[2])) / det
	l3 := 1.0 - l1 - l2
	return l1*p1[1] + l2*p2[1] + l3*p3[1]
}

func (t *Terrain) setNormals(data []float32) []float32 {
	for x := 0; x < VertexCount-1; x++ {
		for z := 0; z < VertexCount-1; z++ {
			heightL := getY(x-1, z, data)
			heightR := getY(x+1, z, data)
			heightD := getY(x, z-1, data)
			heightT := getY(x, z+1, data)

			normal := mgl32.Vec3{heightL - heightR, 2 * t.gridSizeSquare, heightD - heightT}.Normalize()

			normXPos := (x*stride + z*stride*VertexCount) + 3
			data[normXPos] = normal[0]
			data[normXPos+1] = normal[1]
			data[normXPos+2] = normal[2]
		}
	}
	return data
}

func getY(x, z int, data []float32) float32 {
	pos := x*stride + z*stride*VertexCount + 1
	if pos < 0 || pos > len(data) {
		return 0
	}
	return data[pos]
}
