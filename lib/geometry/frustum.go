package geometry

type Frustum [6]Plane

func (f Frustum) Normalise() {
	f[0].Normalise()
	f[1].Normalise()
	f[2].Normalise()
	f[3].Normalise()
	f[4].Normalise()
	f[5].Normalise()
}

func (f Frustum) AABBInside(aabb *AABB) bool {
	for i := 0; i < 6; i++ {
		// compute distance of box center from plane
		d := f[i].Dot(aabb.C())
		// compute the projection interval radius of b onto plane
		r := aabb.R()[0]*abs(f[i][0]) + aabb.R()[1]*abs(f[i][1]) + aabb.R()[2]*abs(f[i][2])
		if d+r <= -f[i][3] {
			return false
		}
	}
	return true
}

func (f Frustum) PointInside(pt [3]float32) bool {
	var d float32
	for i := 0; i < 6; i++ {
		f[i].DistanceToPoint(pt, &d)
		if d <= 0 {
			return false
		}
	}
	return true
}

func (f Frustum) SphereInside(pt [3]float32, radius float32) bool {
	var d float32
	for i := 0; i < 6; i++ {
		f[i].DistanceToPoint(pt, &d)
		if d <= -radius {
			return false
		}
	}
	return true
}

func abs(x float32) float32 {
	if x > 0 {
		return x
	}
	return -x
}
