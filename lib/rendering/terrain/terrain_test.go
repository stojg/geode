package terrain

import (
	"testing"
)

func BenchmarkGenerateTerrain(b *testing.B) {
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		New(0, 0)
	}
}

var res float32

func BenchmarkTerrain_Height(b *testing.B) {
	t := New(0, 0)
	b.ReportAllocs()
	b.ResetTimer()
	var r float32
	for n := 0; n < b.N; n++ {
		r = t.Height(10, 10)
	}
	res = r
}
