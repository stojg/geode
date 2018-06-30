package terrain

import (
	"math"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/stojg/graphics/lib/rendering"
)

const Size float32 = 1024
const VertexCount = 128

func New(gridX, gridZ float32) *Terrain {
	t := &Terrain{
		worldX:         gridX * Size,
		worldZ:         gridZ * Size,
		gridSizeSquare: Size / (float32(VertexCount) - 1),
	}

	t.v, t.i = t.generateTerrain(gridX, gridZ)
	return t
}

type Terrain struct {
	worldX, worldZ float32
	heights        [VertexCount][VertexCount]float32
	gridSizeSquare float32
	v              []float32
	i              []uint32
	mesh           *rendering.Mesh
	texture        interface{}
}

func (t *Terrain) generateTerrain(gridX, gridZ float32) ([]float32, []uint32) {
	xOffset = float64(gridX * float32(VertexCount-1))
	zOffset = float64(gridZ * float32(VertexCount-1))

	hg := NewHeightGenerator(22)

	// https://www.youtube.com/watch?v=yNYwZMmgTJk&list=PLRIWtICgwaX0u7Rf9zkZhLoLuZVfUksDP&index=14
	// @todo calculate size
	var vertices []float32
	for i := 0; i < VertexCount; i++ {
		for j := 0; j < VertexCount; j++ {
			height := hg.Height(float64(i)+xOffset, float64(j)+zOffset)
			t.heights[j][i] = height
			vertices = append(vertices, float32(j)/(float32(VertexCount)-1)*Size)
			vertices = append(vertices, float32(height))
			vertices = append(vertices, float32(i)/(float32(VertexCount)-1)*Size)
			vertices = append(vertices, -66)
			vertices = append(vertices, -66)
			vertices = append(vertices, -66)
			vertices = append(vertices, float32(j)/(float32(VertexCount)-1))
			vertices = append(vertices, float32(i)/(float32(VertexCount)-1))
		}
	}

	vertices = setNormals(vertices)

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

func (t *Terrain) Height(worldX, worldZ float32) float32 {
	terrainX := worldX - t.worldX
	terrainZ := worldZ - t.worldZ

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

func (t *Terrain) Mesh() *rendering.Mesh {
	if t.mesh == nil {
		mesh := rendering.NewMesh()
		mesh.SetVertices(rendering.ConvertToVertices(t.v, t.i), t.i)
		t.mesh = mesh
	}
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

func setNormals(data []float32) []float32 {
	for x := 0; x < VertexCount-1; x++ {
		for z := 0; z < VertexCount-1; z++ {
			heightL := getY(x-1, z, data)
			heightR := getY(x+1, z, data)
			heightD := getY(x, z-1, data)
			heightT := getY(x, z+1, data)
			normal := mgl32.Vec3{heightL - heightR, 2, heightD - heightT}
			normal.Normalize()
			normXPos := (x*8 + z*8*VertexCount) + 3
			data[normXPos] = normal[0]
			data[normXPos+1] = normal[1]
			data[normXPos+2] = normal[2]
		}
	}
	return data
}

func getY(x, z int, data []float32) float32 {
	pos := x*8 + z*8*(VertexCount) + 1
	if pos < 0 || pos > len(data) {
		return 0
	}
	return data[pos]
}
