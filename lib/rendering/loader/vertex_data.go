package loader

import "strconv"

type vertexData struct {
	Type         dataType
	Declarations []*declaration

	meta map[dataType]string
}

func (f *vertexData) setMeta(t dataType, value string) {
	if f.meta == nil {
		f.meta = make(map[dataType]string)
	}
	f.meta[t] = value
}

func (f *vertexData) getMeta(t dataType) string {
	if f.meta != nil {
		return f.meta[t]
	}
	return ""
}

func (f *vertexData) index(index int) *declaration {
	if index >= 0 && index <= 3 {
		for index >= len(f.Declarations) {
			f.Declarations = append(f.Declarations, &declaration{})
		}
		return f.Declarations[index]
	}
	return nil
}

func (vt *vertexData) string() (out string) {

	switch vt.Type {

	case tLine, tPoint:
		hasUVs := false
		if vt.Type == tLine {
			for _, decl := range vt.Declarations {
				if decl.index(tUV) != 0 {
					hasUVs = true
					break
				}
			}
		}
		var prev *declaration
		for di, decl := range vt.Declarations {
			// remove consecutive duplicate points eg. "l 1 1 2 2 3 4 4"
			if prev != nil && prev.equals(decl) {
				continue
			}
			if di > 0 {
				out += " "
			}
			out += strconv.Itoa(decl.index(tVertext))
			if hasUVs {
				out += "/"
				if index := decl.index(tUV); index != 0 {
					out += strconv.Itoa(index)
				}
			}
			prev = decl
		}

	case tFace:
		hasUVs, hasNormals := false, false

		// always use ptr refs if available.
		// this enables simple index rewrites.
		for _, decl := range vt.Declarations {
			if !hasUVs {
				hasUVs = decl.index(tUV) != 0
			}
			if !hasNormals {
				hasNormals = decl.index(tNormal) != 0
			}
			if hasUVs && hasNormals {
				break
			}
		}
		for di, decl := range vt.Declarations {
			if di > 0 {
				out += " "
			}
			out += strconv.Itoa(decl.index(tVertext))
			if hasUVs || hasNormals {
				out += "/"
				if index := decl.index(tUV); index != 0 {
					out += strconv.Itoa(index)
				}
			}
			if hasNormals {
				out += "/"
				if index := decl.index(tNormal); index != 0 {
					out += strconv.Itoa(index)
				}
			}
		}
	}
	return out
}
