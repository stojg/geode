package loader

// declaration

// zero value means it was not declared, should not be written
// @note exception: if sibling declares it, must be written eg. 1//2
type declaration struct {
	vertex int
	uv     int
	normal int

	// Pointers to actual geometry values.
	// When serialized to String, the index is read from ref
	// if available. This enables easy geometry rewrites.
	refVertex, refUV, refNormal *geometryValue
}
