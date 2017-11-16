package loader

import (
	"math"
	"strconv"
)

type geometryValue struct {
	Index      int
	Discard    bool
	X, Y, Z, W float64
}

func equals(a, b, epsilon float64) bool {
	return (math.Abs(a-b) <= epsilon)
}

func (gv *geometryValue) string(t dataType) (out string) {
	switch t {
	case tUV:
		out = strconv.FormatFloat(gv.X, 'g', -1, 64) + " " + strconv.FormatFloat(gv.Y, 'g', -1, 64)
	default:
		out = strconv.FormatFloat(gv.X, 'g', -1, 64) + " " + strconv.FormatFloat(gv.Y, 'g', -1, 64) + " " + strconv.FormatFloat(gv.Z, 'g', -1, 64)
	}
	// omit default values
	switch t {
	case tVertext, tPoint:
		if !equals(gv.W, 1, 1e-10) {
			out += " " + strconv.FormatFloat(gv.W, 'g', -1, 64)
		}
	}
	return out
}

func (gv *geometryValue) distance(to *geometryValue) float64 {
	dx := gv.X - to.X
	dy := gv.Y - to.Y
	dz := gv.Z - to.Z
	return dx*dx + dy*dy + dz*dz
}

func (gv *geometryValue) equals(other *geometryValue, epsilon float64) bool {
	if math.Abs(gv.X-other.X) <= epsilon &&
		math.Abs(gv.Y-other.Y) <= epsilon &&
		math.Abs(gv.Z-other.Z) <= epsilon &&
		math.Abs(gv.W-other.W) <= epsilon {
		return true
	}
	return false
}
