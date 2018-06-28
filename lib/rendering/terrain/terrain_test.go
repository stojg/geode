package terrain

import "testing"

// https://medium.com/@hackintoshrao/daily-code-optimization-using-benchmarks-and-profiling-in-golang-gophercon-india-2016-talk-874c8b4dc3c5

//BenchmarkNew-8   	50000000	        25.5 ns/op
func BenchmarkGetNoise(b *testing.B) {
	for n := 0; n < b.N; n++ {
		getNoise(0, 0)
	}
}

func BenchmarkInterpolate(b *testing.B) {
	for n := 0; n < b.N; n++ {
		interpolate(2.4, 3.5, 0.5)
	}
}

func BenchmarkGenerateTerrain(b *testing.B) {
	for n := 0; n < b.N; n++ {
		generateTerrain(0, 0)
	}
}
