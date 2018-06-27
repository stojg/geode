package terrain

import (
	"math"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/ojrac/opensimplex-go"
	"github.com/stojg/graphics/lib/rendering"
)

const Size float32 = 800
const VertexCount int = 128

const AMPLITUDE float64 = 20
const OCTAVES = 3
const ROUGHNESS = 0.3

var noiseGen *opensimplex.Noise

func init() {
	noiseGen = opensimplex.NewWithSeed(1)

}

func New(gridX, gridZ float32) *Terrain {
	t := &Terrain{
		x: gridX * Size,
		z: gridZ * Size,
	}

	v, i := generateTerrain(gridX, gridZ)
	mesh := rendering.NewMesh()
	mesh.SetVertices(rendering.ConvertToVertices(v, i), i)

	t.mesh = mesh
	return t
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

var xOffset float64 = 0
var zOffset float64 = 0
var seed float64 = 12

func generateTerrain(gridX, gridZ float32) ([]float32, []uint32) {
	xOffset = float64(gridX * float32(VertexCount-1))
	zOffset = float64(gridZ * float32(VertexCount-1))

	// https://www.youtube.com/watch?v=yNYwZMmgTJk&list=PLRIWtICgwaX0u7Rf9zkZhLoLuZVfUksDP&index=14
	var vertices []float32
	for i := 0; i < VertexCount; i++ {
		for j := 0; j < VertexCount; j++ {
			vertices = append(vertices, float32(j)/(float32(VertexCount)-1)*Size)
			vertices = append(vertices, generateHeight(float64(j), float64(i)))
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
	pos := x*8 + z*8*(VertexCount)
	if pos < 0 || pos > len(data) {
		return 0
	}
	return data[pos+1]
}

func generateHeight(x, z float64) float32 {
	total := 0.0
	d := math.Pow(2, OCTAVES-1)
	for i := 0.0; i < OCTAVES; i++ {
		freq := math.Pow(2, i) / d
		amp := math.Pow(ROUGHNESS, float64(i)) * AMPLITUDE
		total += getInterpolatedNoise((x+xOffset)*freq, (z+zOffset)*freq) * amp
	}
	return float32(total)
}

func getInterpolatedNoise(x, z float64) float64 {
	intX, fracX := math.Modf(x)
	intZ, fracZ := math.Modf(z)
	v1 := getSmoothNoise(intX, intZ)
	v2 := getSmoothNoise(intX+1, intZ)
	v3 := getSmoothNoise(intX, intZ+1)
	v4 := getSmoothNoise(intX+1, intZ+1)
	i1 := interpolate(v1, v2, fracX)
	i2 := interpolate(v3, v4, fracX)
	return interpolate(i1, i2, fracZ)
}

func interpolate(a, b, blend float64) float64 {
	theta := blend * math.Pi
	f := (1 - math.Cos(theta)) * 0.5
	return a*(1-f) + b*f
}

func getSmoothNoise(x, z float64) float64 {
	corners := (getNoise(x-1, z-1) + getNoise(x+1, z-1) + getNoise(x-1, z+1) + getNoise(x+1, z+1)) / 16
	sides := (getNoise(x-1, z) + getNoise(x+1, z) + getNoise(x, z-1) + getNoise(x, z+1)) / 8
	center := getNoise(x, z) / 4
	return corners + sides + center
}

func getNoise(x, z float64) float64 {

	return noiseGen.Eval2(x, z)
	//return rand.Float64()*2 - 1
}
