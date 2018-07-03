package resources

func CrossesPlane(plane [3][3]float32, in [3]float32) bool {
	normal := getNormal(plane)
	D := dot(minus(normal), plane[0])
	dist := dot(normal, in) + D
	return dist >= 0
}

func getNormal(a [3][3]float32) [3]float32 {
	edge0 := rawEdge(a[1], a[0])
	edge1 := rawEdge(a[2], a[0])
	return normalise(cross(edge0, edge1))
}

func cross(a, b [3]float32) [3]float32 {
	return [3]float32{a[1]*b[2] - a[2]*b[1], a[2]*b[0] - a[0]*b[2], a[0]*b[1] - a[1]*b[0]}
}

func rawEdge(a, b [3]float32) [3]float32 {
	return [3]float32{a[0] - b[0], a[1] - b[1], a[2] - b[2]}
}

func dot(a, b [3]float32) float32 {
	return a[0]*b[0] + a[1]*b[1] + a[2]*b[2]
}

func minus(a [3]float32) [3]float32 {
	return [3]float32{-a[0], -a[1], -a[2]}
}

// http://gamedevs.org/uploads/fast-extraction-viewing-frustum-planes-from-world-view-projection-matrix.pdf
//func extract_planes_from_projmat(mat[4][4]float, float left[4], float right[4], float top[4], float bottom[4], float near[4], float far[4]) {
//for (int i = 4; i--; ) left[i]      = mat[i][3] + mat[i][0];
//for (int i = 4; i--; ) right[i]     = mat[i][3] - mat[i][0];
//for (int i = 4; i--; ) bottom[i]    = mat[i][3] + mat[i][1];
//for (int i = 4; i--; ) top[i]       = mat[i][3] - mat[i][1];
//for (int i = 4; i--; ) near[i]      = mat[i][3] + mat[i][2];
//for (int i = 4; i--; ) far[i]       = mat[i][3] - mat[i][2];
//}
