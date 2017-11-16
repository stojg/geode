package loader

import (
	"fmt"
	"strconv"
	"strings"
)

type object struct {
	Type       dataType
	Name       string
	Material   string
	VertexData []*vertexData
	Comments   []string

	parent *objectFile
}

// Reads a vertex data line eg. f and l into this object.
//
// If parent objectFile is non nil, additionally converts negative index
// references into absolute indexes and check out of bounds errors.
func (o *object) readVertexData(t dataType, value string, strict bool) (*vertexData, error) {
	var (
		vt  *vertexData
		err error
	)
	switch t {
	case tFace:
		vt, err = parseFaceVertexData(value, strict)
	case tLine, tPoint:
		vt, err = parseListVertexData(t, value, strict)
	default:
		err = fmt.Errorf("object.readVertexData: unsupported vertex data declaration %s %s", t, value)
	}

	if err != nil {
		return nil, err
	} else if o.parent == nil {
		return vt, nil
	}

	// objectFile index references start from 1 not zero.
	// Negative values are relative from the end of currently declared geometry. Convert relative values to absolute.
	geomStats := o.parent.Geometry.stats()
	for _, decl := range vt.Declarations {

		if decl.vertex != 0 {
			if decl.vertex < 0 {
				decl.vertex = decl.vertex + geomStats.Vertices + 1
			}
			if decl.vertex <= 0 || decl.vertex > geomStats.Vertices {
				return nil, fmt.Errorf("vertex index %d out of bounds, %d declared so far", decl.vertex, geomStats.Vertices)
			}
			decl.refVertex = o.parent.Geometry.Vertices[decl.vertex-1]
			if decl.refVertex.Index != decl.vertex {
				return nil, fmt.Errorf("vertex index %d does not match with referenced geometry value %#v", decl.vertex, decl.refVertex)
			}
		}

		if decl.uv != 0 {
			if decl.uv < 0 {
				decl.uv = decl.uv + geomStats.UVs + 1
			}
			if decl.uv <= 0 || decl.uv > geomStats.UVs {
				return nil, fmt.Errorf("uv index %d out of bounds, %d declared so far", decl.uv, geomStats.UVs)
			}
			decl.refUV = o.parent.Geometry.UVs[decl.uv-1]
			if decl.refUV.Index != decl.uv {
				return nil, fmt.Errorf("uv index %d does not match with referenced geometry value %#v", decl.uv, decl.refUV)
			}
		}

		if decl.normal != 0 {
			if decl.normal < 0 {
				decl.normal = decl.normal + geomStats.Normals + 1
			}
			if decl.normal <= 0 || decl.normal > geomStats.Normals {
				return nil, fmt.Errorf("normal index %d out of bounds, %d declared so far", decl.normal, geomStats.Normals)
			}
			decl.refNormal = o.parent.Geometry.Normals[decl.normal-1]
			if decl.refNormal.Index != decl.normal {
				return nil, fmt.Errorf("normal index %d does not match with referenced geometry value %#v", decl.normal, decl.refNormal)
			}
		}
	}
	o.VertexData = append(o.VertexData, vt)
	return vt, nil
}

func parseFaceVertexData(str string, strict bool) (vt *vertexData, err error) {
	vt = &vertexData{
		Type: tFace,
	}
	for iMain, part := range strings.Split(str, " ") {
		dest := vt.index(iMain)
		if dest == nil {
			if strict {
				return nil, fmt.Errorf("Invalid face index %d in %s", iMain, str)
			}
			break
		}
		for iPart, datapart := range strings.Split(part, "/") {
			value := 0
			// can be empty eg. "f 1//1 2//2 3//3 4//4"
			if len(datapart) > 0 {
				value, err = strconv.Atoi(datapart)
				if err != nil {
					return nil, err
				}
			}
			switch iPart {
			case 0:
				dest.vertex = value
			case 1:
				dest.uv = value
			case 2:
				dest.normal = value
			default:
				if strict {
					return nil, fmt.Errorf("Invalid face vertex data index %d.%d in %s", iMain, iPart, str)
				}
				break
			}
		}
	}
	return vt, nil
}

func parseListVertexData(t dataType, str string, strict bool) (*vertexData, error) {
	if t != tLine && t != tPoint {
		return nil, fmt.Errorf("parseListVertexData supports face and point type, given: %s", t)
	}
	vt := &vertexData{
		Type: t,
	}
	for iMain, part := range strings.Split(str, " ") {
		decl := &declaration{}
		for iPart, datapart := range strings.Split(part, "/") {
			if len(datapart) == 0 {
				continue
			}
			value, vErr := strconv.Atoi(datapart)
			if vErr != nil {
				return nil, vErr
			}
			switch iPart {
			case 0:
				decl.vertex = value
			case 1:
				decl.uv = value
			default:
				if strict {
					return nil, fmt.Errorf("Invalid face vertex data index %d.%d in %s", iMain, iPart, str)
				}
				break
			}
		}
		vt.Declarations = append(vt.Declarations, decl)
	}
	return vt, nil
}
