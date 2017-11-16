package loader

import "fmt"

// declaration

// zero value means it was not declared, should not be written
// @note exception: if sibling declares it, must be written eg. 1//2
type declaration struct {
	vertex int
	uv     int
	normal int

	// Pointers to actual geometry values.
	// When serialized to string, the index is read from ref
	// if available. This enables easy geometry rewrites.
	refVertex, refUV, refNormal *geometryValue
}

func (d *declaration) equals(other *declaration) bool {
	if d.index(tVertext) != other.index(tVertext) ||
		d.index(tUV) != other.index(tUV) ||
		d.index(tNormal) != other.index(tNormal) {
		return false
	}
	return true
}

// Use this getter when possible index rewrites has occurred.
// Will first return index from geometry value pointers, if available.
func (d *declaration) index(t dataType) int {
	switch t {
	case tVertext:
		if d.refVertex != nil {
			return d.refVertex.Index
		}
		return d.vertex
	case tUV:
		if d.refUV != nil {
			return d.refUV.Index
		}
		return d.uv
	case tNormal:
		if d.refNormal != nil {
			return d.refNormal.Index
		}
		return d.normal
	default:
		fmt.Printf("declaration.index: unsupported type %s\n", t)
	}
	return 0
}
