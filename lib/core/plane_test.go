package core

import (
	"testing"

	"github.com/go-gl/mathgl/mgl32"
)

func TestDistanceToPoint(t *testing.T) {

	type pos [3]float32
	var plane [4]float32

	plane = [4]float32{0, 0, -1, 1}
	point := pos{0, 0, 1}
	var d float32
	dotPlane(plane, point, &d)
	var inside bool
	if d >= 0 {
		inside = true
	}
	if inside != true {
		t.Errorf("expected %t, got %t", true, inside)
	}

	point = pos{0, 0, -10}
	dotPlane(plane, point, &d)
	if d >= 0 {
		inside = false
	}
	if inside != false {
		t.Errorf("expected %t, got %t", false, inside)
	}

}

var d float32

func BenchmarkDotPlane(b *testing.B) {
	plane := [4]float32{3, 0, 0, 9}
	pt := [3]float32{1, 1, 2}

	b.ReportAllocs()
	b.ResetTimer()
	var x float32
	for i := 0; i < b.N; i++ {
		var x float32
		dotPlane(plane, pt, &x)
	}
	d = x
}

func BenchmarkIsVisible(b *testing.B) {
	planes := [6][4]float32{
		{-0.9870641, 0.048430424, 0.15283628, 11.399004},
		{0.40589866, 0.048430473, 0.9126339, 5.067353},
		{-0.26035964, 0.8392691, 0.47732612, 7.3768578},
		{-0.32080588, -0.74240816, 0.5881443, 9.089501},
		{-0.4773343, 0.079555705, 0.8751131, 13.424475},
		{0.4773718, -0.07955561, -0.8750926, 498.57483},
	}

	aabb := [3][2]float32{
		{256, 256}, {10.220467, -0.55961514}, {256, 256},
	}

	transform := mgl32.Mat4{
		1.000000, 0.000000, 0.000000, -256.000000,
		0.000000, 1.000000, 0.000000, 0.000000,
		0.000000, 0.000000, 1.000000, -256.000000,
		0.000000, 0.000000, 0.000000, 1.000000,
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		IsVisible(planes, aabb, transform)
	}
}
