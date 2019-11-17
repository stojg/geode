package math

import (
	"testing"

	"github.com/go-gl/mathgl/mgl32"
)

var R mgl32.Mat4

func BenchmarkMul4(b *testing.B) {
	matA := mgl32.Ident4()
	matB := mgl32.Ident4()

	var r mgl32.Mat4

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		Mul4(matA, matB, &r)
	}

	R = r
}
