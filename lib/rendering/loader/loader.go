// package loader loads .obj files
//
// http://www.martinreddy.net/gfx/3d/OBJ.spec
// https://github.com/jonnenauha/obj-simplify/blob/master/objectfile/structs.go
package loader

import (
	"fmt"
	"os"
)

func Load(filename string) ([]float32, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	obj, num, err := parse(f)

	if err != nil {
		return nil, fmt.Errorf("error in '%s' at line: %d", filename, num)
	}

	for _, object := range obj.Objects {
		var data []float32
		// convert the face data into actual data ready for openGL loading
		for _, vert := range object.VertexData {
			data = add(data, vert.Declarations[0])
			data = add(data, vert.Declarations[1])
			data = add(data, vert.Declarations[2])
			for i := 3; i < len(vert.Declarations); i++ {
				data = add(data, vert.Declarations[i-3])
				data = add(data, vert.Declarations[i-1])
				data = add(data, vert.Declarations[i])
			}
		}
		return data, nil
	}
	return nil, nil
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
	return []float32{float32(val.X), float32(val.Y), float32(val.Z), float32(val.Z)}
}
