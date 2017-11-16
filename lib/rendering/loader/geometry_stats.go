package loader

type geometryStats struct {
	Vertices, Normals, UVs, Params int
}

func (gs geometryStats) isEmpty() bool {
	return gs.Vertices == 0 && gs.UVs == 0 && gs.Normals == 0 && gs.Params == 0
}

func (gs geometryStats) num(t dataType) int {
	switch t {
	case tVertext:
		return gs.Vertices
	case tUV:
		return gs.UVs
	case tNormal:
		return gs.Normals
	case tParam:
		return gs.Params
	default:
		return 0
	}
}
