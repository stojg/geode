package resources

import (
	"reflect"
	"testing"
)

func TestConvertToVertices(t *testing.T) {
	type args struct {
		meshdata []float32
		indices  []uint32
	}
	tests := []struct {
		name string
		args args
		want []Vertex
	}{
		{
			name: "basic",
			args: args{
				meshdata: []float32{
					0, 0, 0, 1, 1, 1, 0, 1,
					1, 0, 0, 1, 1, 1, 1, 0,
					0, 1, 0, 1, 1, 1, 1, 1,
				},
				indices: []uint32{0, 1, 2},
			},
			want: []Vertex{
				{Pos: [3]float32{0, 0, 0}, Normal: [3]float32{1, 1, 1}, TexCoords: [2]float32{0, 1}, Tangent: [3]float32{0, 1, 0}},
				{Pos: [3]float32{1, 0, 0}, Normal: [3]float32{1, 1, 1}, TexCoords: [2]float32{1, 0}, Tangent: [3]float32{0, 1, 0}},
				{Pos: [3]float32{0, 1, 0}, Normal: [3]float32{1, 1, 1}, TexCoords: [2]float32{1, 1}, Tangent: [3]float32{0, 1, 0}},
			},
		},
		{
			name: "reuse",
			args: args{
				meshdata: []float32{
					0, 0, 0, 1, 1, 1, 0, 1,
					1, 0, 0, 1, 1, 1, 1, 0,
					0, 1, 0, 1, 1, 1, 1, 1,
					1, 1, 0, 1, 1, 1, 1, 1,
				},
				indices: []uint32{0, 1, 2, 0, 1, 3},
			},
			want: []Vertex{
				{Pos: [3]float32{0, 0, 0}, Normal: [3]float32{1, 1, 1}, TexCoords: [2]float32{0, 1}, Tangent: [3]float32{0, 1, 0}},
				{Pos: [3]float32{1, 0, 0}, Normal: [3]float32{1, 1, 1}, TexCoords: [2]float32{1, 0}, Tangent: [3]float32{0, 1, 0}},
				{Pos: [3]float32{0, 1, 0}, Normal: [3]float32{1, 1, 1}, TexCoords: [2]float32{1, 1}, Tangent: [3]float32{0, 1, 0}},
				{Pos: [3]float32{1, 1, 0}, Normal: [3]float32{1, 1, 1}, TexCoords: [2]float32{1, 1}, Tangent: [3]float32{0, 0, 0}},
			},
		},
		{
			name: "skewed",
			args: args{
				meshdata: []float32{
					0, 0, 0, 1, 1, 1, 1, 0,
					1, 0, 0, 1, 1, 1, 0, 1,
					0, 1, 0, 1, 1, 1, 1, 1,
				},
				indices: []uint32{0, 2, 1},
			},
			want: []Vertex{
				{Pos: [3]float32{0, 0, 0}, Normal: [3]float32{1, 1, 1}, TexCoords: [2]float32{1, 0}, Tangent: [3]float32{-0.70710677, 0.70710677, 0}},
				{Pos: [3]float32{1, 0, 0}, Normal: [3]float32{1, 1, 1}, TexCoords: [2]float32{0, 1}, Tangent: [3]float32{-0.70710677, 0.70710677, 0}},
				{Pos: [3]float32{0, 1, 0}, Normal: [3]float32{1, 1, 1}, TexCoords: [2]float32{1, 1}, Tangent: [3]float32{-0.70710677, 0.70710677, 0}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ConvertToVertices(tt.args.meshdata, tt.args.indices)
			if len(got) != len(tt.want) {
				t.Errorf("ConvertToVertices() want %d vertices, got %d", len(got), len(tt.want))
			} else if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConvertToVertices() failed")
				for i := range got {
					t.Errorf("got %+v, want %+v", got[i], tt.want[i])
				}
			}
		})
	}
}

var tgot []Vertex

func BenchmarkConvertToVertices(b *testing.B) {
	b.ReportAllocs()
	meshdata := []float32{
		0, 0, 0, 1, 1, 1, 1, 0,
		1, 0, 0, 1, 1, 1, 0, 1,
		0, 1, 0, 1, 1, 1, 1, 1,
	}
	indices := []uint32{0, 2, 1}

	var got []Vertex
	for i := 0; i < b.N; i++ {
		got = ConvertToVertices(meshdata, indices)
	}
	tgot = got
}
