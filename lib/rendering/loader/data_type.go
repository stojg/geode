package loader

type dataType int

const (
	tUnkown dataType = iota

	tComment        // #
	tMtlLib         // mtllib
	tMtlUse         // usemtl
	tChildGroup     // g
	tChildObject    // o
	tSmoothingGroup // s
	tVertext        // v
	tNormal         // vn
	tUV             // vt
	tParam          // vp
	tFace           // f
	tLine           // l
	tPoint          // p
	tCurve          // curv
	tCurve2         // curv2
	tSurface        // surf
)

func (ot dataType) string() string {
	switch ot {
	case tComment:
		return "#"
	case tMtlLib:
		return "mtllib"
	case tMtlUse:
		return "usemtl"
	case tChildGroup:
		return "g"
	case tChildObject:
		return "o"
	case tSmoothingGroup:
		return "s"
	case tVertext:
		return "v"
	case tNormal:
		return "vn"
	case tUV:
		return "vt"
	case tParam:
		return "vp"
	case tFace:
		return "f"
	case tLine:
		return "l"
	case tPoint:
		return "p"
	case tCurve:
		return "curv"
	case tCurve2:
		return "curv2"
	case tSurface:
		return "surf"
	}
	return ""
}

func (ot dataType) name() string {
	switch ot {
	case tVertext:
		return "vertices"
	case tNormal:
		return "normals"
	case tUV:
		return "uvs"
	case tParam:
		return "params"
	case tChildGroup:
		return "group"
	case tChildObject:
		return "object"
	}
	return ""
}

func typeFromString(str string) dataType {
	switch str {
	case "#":
		return tComment
	case "mtllib":
		return tMtlLib
	case "usemtl":
		return tMtlUse
	case "g":
		return tChildGroup
	case "o":
		return tChildObject
	case "s":
		return tSmoothingGroup
	case "v":
		return tVertext
	case "vn":
		return tNormal
	case "vt":
		return tUV
	case "vp":
		return tParam
	case "f":
		return tFace
	case "l":
		return tLine
	case "p":
		return tPoint
	case "curv":
		return tCurve
	case "curv2":
		return tCurve2
	case "surf":
		return tSurface
	}
	return tUnkown
}
