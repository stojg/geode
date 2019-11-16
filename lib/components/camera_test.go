package components

import (
	"testing"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/stojg/geode/lib/geometry"
)

var lres geometry.Frustum

func BenchmarkExtractPlanesFromProjection(b *testing.B) {
	var l geometry.Frustum
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
