package images

import (
	"image"
	"testing"
)

func TestFlip(t *testing.T) {
	src := []uint8{
		1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, // width 3
		13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24,
	}
	rgba := image.NewRGBA(image.Rect(0, 0, 3, 2))
	rgba.Pix = src
	Flip(rgba)
	expected := []uint8{
		13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24,
		1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12,
	}
	if len(rgba.Pix) != len(src) {
		t.Fatalf("expected same size result, src %d, result %d", len(src), len(rgba.Pix))
	}
	for i := range expected {
		if rgba.Pix[i] != expected[i] {
			t.Fatalf("Expected:\n%v\ngot:\n%v\nposition %d", expected, rgba.Pix, i)
		}
	}
}

// BenchmarkFlip-8   	   10000	    196042 ns/op	    4096 B/op	       1 allocs/op - scuffed-plastic-ao.png
// BenchmarkFlip-8   	30000000	        35.5 ns/op	      16 B/op	       1 allocs/op - small
func BenchmarkFlip(b *testing.B) {
	src := []uint8{
		1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, // width 3
		13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24,
	}
	rgba := image.NewRGBA(image.Rect(0, 0, 3, 2))
	rgba.Pix = src

	//rgba, err := RGBAImagedata("../../res/textures/scuffed-plastic/scuffed-plastic-ao.png")
	//if err != nil {
	//	b.Fatal(err)
	//}

	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		Flip(rgba)
	}
}

func TestFlipRaw(t *testing.T) {
	src := []float32{
		1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, // width 4
		13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, //
	}
	FlipRaw(4, 2, src)
	expected := []float32{
		13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24,
		1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12,
	}
	for i := range src {
		if src[i] != expected[i] {
			t.Fatalf("Expected:\n%v\ngot:\n%v\nposition %d", expected, src, i)
		}
	}
}

// BenchmarkFlipRaw-8   	     100	  10292038 ns/op	   49152 B/op	       1 allocs/op - sky0016.hdr
// BenchmarkFlipRaw-8   	30000000	        43.2 ns/op	      48 B/op	       1 allocs/op - small
func BenchmarkFlipRaw(b *testing.B) {
	//w, h, in, err := RGBEImagedata("../../res/textures/sky0016.hdr")
	//if err != nil {
	//	b.Fatal(err)
	//}
	src := []float32{
		1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, // width 4
		13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, //
	}
	w, h := 4, 2
	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		FlipRaw(w, h, src)
	}
}
