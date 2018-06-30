// package loader loads .obj files
//
// http://www.martinreddy.net/gfx/3d/OBJ.spec
// https://github.com/jonnenauha/obj-simplify/blob/master/objectfile/structs.go
package loader

import (
	"fmt"
	"os"
)

// Load takes in a wavefront .obj filename  and translates it into data usable in opengl, it returns a slice of slice
// of unique vertices and a slice of slice of indices that points to that data.
// An .obj file can contain multiple objects (i.e. meshes) and hence the need for a slice of slice
// The vertices is a continuous list of float32 in the order of position (x,y,z), normals (x,y,z) and texture uv
// coordinates (u,v)
func Load(filename string) ([][]float32, [][]uint32, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, nil, err
	}
	defer f.Close()

	obj, num, err := parse(f)

	if err != nil {
		return nil, nil, fmt.Errorf("error in '%s' at line: %d", filename, num)
	}

	var perObjectVertices [][]float32
	var perObjectIndices [][]uint32

	for _, object := range obj.Objects {
		indices, uniqueVertices := calculateIndices(object)
		perObjectVertices = append(perObjectVertices, uniqueVertices)
		perObjectIndices = append(perObjectIndices, indices)
	}
	return perObjectVertices, perObjectIndices, nil
}

func calculateIndices(object *object) ([]uint32, []float32) {
	// this is used as a lookup map to check if a vertex has been indexed
	existing := make(map[[3]int]uint32)
	var indices []uint32
	var uniqueVertices []float32
	// triangularize the wavefront .obj file data and calculate indices to reduce vertices
	for _, v := range object.VertexData {
		// first triangle in the face/n-gon
		for i := 0; i < 3; i++ {
			uniqueVertices, indices = getUniqueVertices(v.Declarations[i], existing, indices, uniqueVertices)
		}
		// calculate the other triangles in the face/n-gon
		for i := 3; i < len(v.Declarations); i++ {
			uniqueVertices, indices = getUniqueVertices(v.Declarations[i-3], existing, indices, uniqueVertices)
			uniqueVertices, indices = getUniqueVertices(v.Declarations[i-1], existing, indices, uniqueVertices)
			uniqueVertices, indices = getUniqueVertices(v.Declarations[i], existing, indices, uniqueVertices)
		}
	}
	return indices, uniqueVertices
}

func getUniqueVertices(decl *declaration, lookupMap map[[3]int]uint32, indices []uint32, uniqueVertices []float32) ([]float32, []uint32) {
	key := [3]int{decl.vertex, decl.normal, decl.uv}
	index, exists := lookupMap[key]
	if !exists {
		index = uint32(len(lookupMap))
		lookupMap[key] = index
		uniqueVertices = declarationToFloat32(decl, uniqueVertices)
	}
	indices = append(indices, index)
	return uniqueVertices, indices
}

// adds the position, normal and texture uv values to the slice data and return it
func declarationToFloat32(in *declaration, data []float32) []float32 {
	data = appendGeometry(in.refVertex, data, 3)
	data = appendGeometry(in.refNormal, data, 3)
	if in.refUV != nil {
		data = appendGeometry(in.refUV, data, 2)
	} else {
		data = append(data, 0, 0)
	}
	return data
}

// appendGeometry appends the first `count` values of the `in` to the `data` slice and returns the result
func appendGeometry(in *geometryValue, data []float32, count int) []float32 {
	return append(data, geometryToFloat32(in)[:count]...)
}

// cast values to float32 to work nicely with opengl
func geometryToFloat32(val *geometryValue) []float32 {
	return []float32{float32(val.x), float32(val.y), float32(val.z), float32(val.z)}
}
