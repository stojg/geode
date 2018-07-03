package resources

import (
	"testing"
)

func TestCrossesPlane(t *testing.T) {

	type pos [3]float32
	var plane [3][3]float32

	plane[0] = [3]float32{0, 1, 0}
	plane[1] = [3]float32{0, 0, 0}
	plane[1] = [3]float32{0, 1, 0}

	point := pos{2, 2, 2}
	actual := CrossesPlane(plane, point)
	expected := false
	if actual != expected {
		t.Errorf("expected %t, got %t", expected, actual)
	}

	point = pos{-2, -2, -2}
	actual = CrossesPlane(plane, point)
	expected = false
	if actual != expected {
		t.Errorf("expected %t, got %t", expected, actual)
	}

}
