package geometry

import "testing"

func TestNormalisePlane(t *testing.T) {
	plane := Plane{3, 0, 0, 9}
	plane.Normalise()
	if plane[3] != 3 {
		t.Errorf("normalisePlane failed: %f", plane[3])
	}
}

func TestDistanceToPoint(t *testing.T) {

	type pos [3]float32
	plane := Plane{0, 0, -1, 1}

	point := pos{0, 0, 1}
	var d float32
	plane.DistanceToPoint(point, &d)

	var inside bool
	if d >= 0 {
		inside = true
	}
	if inside != true {
		t.Errorf("expected %t, got %t", true, inside)
	}

	point = pos{0, 0, -10}
	plane.DistanceToPoint(point, &d)
	if d >= 0 {
		inside = false
	}
	if inside != false {
		t.Errorf("expected %t, got %t", false, inside)
	}
}

var d float32

func BenchmarkDotPlane(b *testing.B) {
	plane := Plane{3, 0, 0, 9}
	pt := [3]float32{1, 1, 2}

	b.ReportAllocs()
	b.ResetTimer()
	var x float32
	for i := 0; i < b.N; i++ {
		var x float32
		plane.DistanceToPoint(pt, &x)
	}
	d = x
}
