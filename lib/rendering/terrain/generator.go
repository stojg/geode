package terrain

import (
	"math"

	"github.com/ojrac/opensimplex-go"
)

const OCTAVES = 5
const ROUGHNESS = 0.5
const AMPLITUDE float64 = 60

var xOffset float64
var zOffset float64

func NewHeightGenerator(seed int64) *HeightGenerator {
	return &HeightGenerator{
		noiseGen: opensimplex.NewWithSeed(seed),
	}
}

type HeightGenerator struct {
	noiseGen *opensimplex.Noise
}

func (h *HeightGenerator) Height(x, z float64) float32 {
	total := 0.0
	d := math.Pow(2, OCTAVES-1)
	for i := 0.0; i < OCTAVES; i++ {
		freq := math.Pow(2, i) / d
		amp := math.Pow(ROUGHNESS, float64(i)) * AMPLITUDE
		total += h.getInterpolatedNoise((x+xOffset)*freq, (z+zOffset)*freq) * amp
	}
	return float32(total)
}

func (h *HeightGenerator) getInterpolatedNoise(x, z float64) float64 {
	intX, fracX := math.Modf(x)
	intZ, fracZ := math.Modf(z)
	v1 := h.getSmoothNoise(intX, intZ)
	v2 := h.getSmoothNoise(intX+1, intZ)
	v3 := h.getSmoothNoise(intX, intZ+1)
	v4 := h.getSmoothNoise(intX+1, intZ+1)
	i1 := h.interpolate(v1, v2, fracX)
	i2 := h.interpolate(v3, v4, fracX)
	return h.interpolate(i1, i2, fracZ)
}

func (h *HeightGenerator) getSmoothNoise(x, z float64) float64 {
	corners := (h.getNoise(x-1, z-1) + h.getNoise(x+1, z-1) + h.getNoise(x-1, z+1) + h.getNoise(x+1, z+1)) / 16
	sides := (h.getNoise(x-1, z) + h.getNoise(x+1, z) + h.getNoise(x, z-1) + h.getNoise(x, z+1)) / 8
	center := h.getNoise(x, z) / 4
	return corners + sides + center
}

func (h *HeightGenerator) getNoise(x, z float64) float64 {
	return h.noiseGen.Eval2(x, z)
}

func (h *HeightGenerator) interpolate(a, b, blend float64) float64 {
	theta := blend * math.Pi
	f := (1 - math.Cos(theta)) * 0.5
	return a*(1-f) + b*f
}
