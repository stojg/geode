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
	verts, inds := calculateIndices(obj)

	if len(verts) != 1 {
		t.Fatalf("expected 1 vertice group, got %d ", len(verts))
	}

	if len(inds) != 1 {
		t.Fatalf("expected 1 indice group, got %d ", len(inds))
	}

	stride := 8 // 3 pos, 3 normals, 2 uv

	verticesData := verts[0]
	if len(verticesData)/stride != 24 {
		t.Fatalf("expected 24 vertices in a cube, got %d ", len(verticesData)/stride)
	}

	//fmt.Println(unsafe.Sizeof(verticesData))

	indices := inds[0]
	if len(indices) != 36 {
		t.Fatalf("expected 36 vertices in a cube, got %d ", len(indices))
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
		calculateIndices(obj)
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
