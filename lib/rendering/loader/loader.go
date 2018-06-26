// package loader loads .obj files
//
// http://www.martinreddy.net/gfx/3d/OBJ.spec
// https://github.com/jonnenauha/obj-simplify/blob/master/objectfile/structs.go
package loader

import (
	"fmt"
	"os"
)

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
		vertexCombinations := make(map[[3]int]uint32)
		var indices []uint32
		var nextIndex uint32 = 0
		var vertices []float32
		// convert the face vertices into actual vertices ready for openGL loading
		for _, vert := range object.VertexData {
			indices, vertices = dosomeThing2(vert.Declarations[0], vertexCombinations, &nextIndex, indices, vertices)
			indices, vertices = dosomeThing2(vert.Declarations[1], vertexCombinations, &nextIndex, indices, vertices)
			indices, vertices = dosomeThing2(vert.Declarations[2], vertexCombinations, &nextIndex, indices, vertices)
			for i := 3; i < len(vert.Declarations); i++ {
				indices, vertices = dosomeThing2(vert.Declarations[i-3], vertexCombinations, &nextIndex, indices, vertices)
				indices, vertices = dosomeThing2(vert.Declarations[i-1], vertexCombinations, &nextIndex, indices, vertices)
				indices, vertices = dosomeThing2(vert.Declarations[i], vertexCombinations, &nextIndex, indices, vertices)
			}
		}
		perObjectVertices = append(perObjectVertices, vertices)
		perObjectIndices = append(perObjectIndices, indices)
	}
	return perObjectVertices, perObjectIndices, nil
}

func dosomeThing2(decl *declaration, vertexCombinations map[[3]int]uint32, nextIndex *uint32, indices []uint32, data []float32) ([]uint32, []float32) {
	idx, exists := getSomething(decl, vertexCombinations, nextIndex)
	indices = append(indices, idx)
	if !exists {
		data = add(data, decl)
	}
	return indices, data
}

func getSomething(in *declaration, indexLookup map[[3]int]uint32, nextIndex *uint32) (uint32, bool) {
	key := [3]int{in.vertex, in.normal, in.uv}
	index, ok := indexLookup[key]
	if !ok {
		index = *nextIndex
		indexLookup[key] = index
		*nextIndex++
	}
	return index, ok
}

func add(data []float32, in *declaration) []float32 {
	data = appendValues(data, in.refVertex, 3)
	data = appendValues(data, in.refNormal, 3)
	if in.refUV != nil {
		data = appendValues(data, in.refUV, 2)
	} else {
		data = append(data, 0, 0)
	}
	return data
}

func appendValues(data []float32, in *geometryValue, count int) []float32 {
	return append(data, toFloat32(in)[:count]...)
}

func toFloat32(val *geometryValue) []float32 {
	return []float32{float32(val.x), float32(val.y), float32(val.z), float32(val.z)}
}
