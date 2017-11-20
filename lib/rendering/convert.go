package rendering

import "math"

func ConvertToVertices(meshdata []float32) []Vertex {
	const stride = 8

	if len(meshdata)%stride != 0 {
		panic("the graph data is not a multiple of 8, want [3]Pos, [3]Normals, [2]TexCoords")
	}
	var vertices []Vertex

	for i := 0; i < len(meshdata); i += stride {
		var vertex Vertex
		copy(vertex.Pos[:], meshdata[i:i+3])
		copy(vertex.Normal[:], meshdata[i+3:i+6])
		copy(vertex.TexCoords[:], meshdata[i+6:i+8])
		vertices = append(vertices, vertex)
	}

	// calculate tangents and bi-tangents
	for i := 0; i < len(vertices); i += 3 {
		v0 := vertices[i]
		v1 := vertices[i+1]
		v2 := vertices[i+2]

		edge1 := edge(v1, v0)
		edge2 := edge(v2, v0)

		deltaU1 := v1.TexCoords[0] - v0.TexCoords[0]
		deltaV1 := v1.TexCoords[1] - v0.TexCoords[1]
		deltaU2 := v2.TexCoords[0] - v0.TexCoords[0]
		deltaV2 := v2.TexCoords[1] - v0.TexCoords[1]

		f := 1.0 / (deltaU1*deltaV2 - deltaU2*deltaV1)

		var tangent [3]float32
		tangent[0] = f * (deltaV2*edge1[0] - deltaV1*edge2[0])
		tangent[1] = f * (deltaV2*edge1[1] - deltaV1*edge2[1])
		tangent[2] = f * (deltaV2*edge1[2] - deltaV1*edge2[2])

		tangent = normalise(tangent)

		copy(vertices[i].Tangent[:], tangent[:])
		copy(vertices[i+1].Tangent[:], tangent[:])
		copy(vertices[i+2].Tangent[:], tangent[:])
	}
	return vertices
}

func edge(a, b Vertex) [3]float32 {
	return [3]float32{
		a.Pos[0] - b.Pos[0],
		a.Pos[1] - b.Pos[1],
		a.Pos[2] - b.Pos[2],
	}
}

func normalise(vec [3]float32) [3]float32 {
	l := 1.0 / float32(math.Sqrt(float64(vec[0]*vec[0]+vec[1]*vec[1]+vec[2]*vec[2])))
	return [3]float32{vec[0] * l, vec[1] * l, vec[2] * l}
}
