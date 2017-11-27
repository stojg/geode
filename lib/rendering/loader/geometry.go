package loader

import (
	"fmt"
	"strconv"
	"strings"
)

type geometry struct {
	Vertices []*geometryValue // v    x y z [w]
	Normals  []*geometryValue // vn   i j k
	UVs      []*geometryValue // vt   u [v [w]]
	Params   []*geometryValue // vp   u v [w]
}

func (g *geometry) readValue(t dataType, value string, strict bool) (*geometryValue, error) {
	gv := &geometryValue{}

	// W is according to the spec by default 1, not serialized in String() if not touched
	// @todo, it doesnt look like Point is supported in the switch t loop at the end?
	if t == tVertext || t == tPoint {
		gv.w = 1
	}

	for i, part := range strings.Split(value, " ") {
		if len(part) == 0 {
			continue
		}
		if part == "-0" {
			part = "0"
		} else if strings.Index(part, "-0.") == 0 {
			// "-0.000000" etc.
			if trimmed := strings.TrimRight(part, "0"); trimmed == "-0." {
				part = "0"
			}
		}
		num, err := strconv.ParseFloat(part, 64)
		if err != nil {
			return nil, fmt.Errorf("found invalid number from %q: %s", value, err)
		}

		switch i {
		case 0:
			gv.x = num
		case 1:
			gv.y = num
		case 2:
			gv.z = num
		case 3:
			if strict && t != tVertext {
				return nil, fmt.Errorf("found invalid fourth component: %s %s", t.String(), value)
			}
			gv.w = num
		default:
			if strict {
				return nil, fmt.Errorf("found invalid fifth component: %s %s", t.String(), value)
			}
		}
	}
	// objectFile refs start from 1 not zero
	gv.index = len(g.get(t)) + 1
	switch t {
	case tVertext:
		g.Vertices = append(g.Vertices, gv)
	case tUV:
		g.UVs = append(g.UVs, gv)
	case tNormal:
		g.Normals = append(g.Normals, gv)
	case tParam:
		g.Params = append(g.Params, gv)
	default:
		return nil, fmt.Errorf("unkown geometry value type %d %s", t, t)
	}
	return gv, nil
}

func (g *geometry) get(t dataType) []*geometryValue {
	switch t {
	case tVertext:
		return g.Vertices
	case tNormal:
		return g.Normals
	case tUV:
		return g.UVs
	case tParam:
		return g.Params
	}
	return nil
}

func (g *geometry) stats() geometryStats {
	return geometryStats{
		Vertices: len(g.Vertices),
		Normals:  len(g.Normals),
		UVs:      len(g.UVs),
		Params:   len(g.Params),
	}
}

func newGeometry() *geometry {
	return &geometry{
		Vertices: make([]*geometryValue, 0),
		Normals:  make([]*geometryValue, 0),
		UVs:      make([]*geometryValue, 0),
		Params:   make([]*geometryValue, 0),
	}
}
