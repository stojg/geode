package geometry

import (
	"testing"

	"github.com/go-gl/mathgl/mgl32"
)

func BenchmarkAABBInView(b *testing.B) {
	f := Frustum{
		{-0.9870641, 0.048430424, 0.15283628, 11.399004},
		{0.40589866, 0.048430473, 0.9126339, 5.067353},
		{-0.26035964, 0.8392691, 0.47732612, 7.3768578},
		{-0.32080588, -0.74240816, 0.5881443, 9.089501},
		{-0.4773343, 0.079555705, 0.8751131, 13.424475},
		{0.4773718, -0.07955561, -0.8750926, 498.57483},
	}

	aabb := &AABB{}
	aabb.SetR(mgl32.Vec3{512, 10, 512})
	aabb.SetC(mgl32.Vec3{0, 5, 0})

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		f.AABBInside(aabb)
	}
}
