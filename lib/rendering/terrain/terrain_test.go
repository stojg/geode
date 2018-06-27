package terrain

import "testing"

func BenchmarkNew(b *testing.B) {
	for n := 0; n < b.N; n++ {
		getNoise(0, 0)
	}
}
