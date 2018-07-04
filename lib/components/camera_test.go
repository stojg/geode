package components

import (
	"testing"

	"github.com/go-gl/mathgl/mgl32"
)

func TestNormalisePlane(t *testing.T) {
	plane := [4]float32{3, 0, 0, 9}
	normalisePlane(&plane)
	if plane[3] != 3 {
		t.Errorf("normalisePlane failed: %f", plane[3])
	}
}

var lres [6][4]float32

func BenchmarkExtractPlanesFromProjection(b *testing.B) {
	var l [6][4]float32
	fov := float32(90)
	width, height := 100, 100
	near := float32(1)
	far := float32(10)
	projection := mgl32.Perspective(mgl32.DegToRad(fov), float32(width/height), near, far)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l = extractPlanesFromProjection(projection, true)
	}

	lres = l
}
