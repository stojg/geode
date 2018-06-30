package loader

import (
	"os"
	"testing"
)

func TestCalculateIndices(t *testing.T) {
	obj, err := getObject("cube.obj")
	if err != nil {
		t.Fatal(err)
	}
	indices, verts := calculateIndices(obj.Objects[0])

	const stride = 8 // 3 pos, 3 normals, 2 uv

	vertices := verts
	if len(vertices)/stride != 24 {
		t.Fatalf("expected 24 vertices in a cube, got %d ", len(vertices)/stride)
	}

	if len(indices) != 36 {
		t.Fatalf("expected 36 indices in a cube, got %d ", len(indices))
	}
}

// BenchmarkCalculateIndices-8   	  200000	      5567 ns/op	    4332 B/op	      19 allocs/op
func BenchmarkCalculateIndices(b *testing.B) {

	obj, err := getObject("cube.obj")
	if err != nil {
		b.Fatal(err)
	}
	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		calculateIndices(obj.Objects[0])
	}
}

func BenchmarkLoad(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		Load("_testdata/sphere.obj")
	}
}

func getObject(name string) (*objectFile, error) {

	f, err := os.Open("_testdata/" + name)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	obj, _, err := parse(f)

	return obj, err
}
