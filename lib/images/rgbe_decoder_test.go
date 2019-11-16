package images

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"flag"
	"io/ioutil"
	"path/filepath"
	"testing"
)

var update = flag.Bool("update", false, "update .golden files")

func Test_readHeader(t *testing.T) {
	testData := `#?RADIANCE
# Made with 100% pure HDR Shop
FORMAT=32-bit_rle_rgbe
EXPOSURE=          1.0000000000000

-Y 2000 +X 4000
`
	r := bytes.NewBufferString(testData)
	reader := bufio.NewReader(r)
	width, height, err := readHeader(reader)
	if err != nil {
		t.Fatal(err)
	}

	if width != 4000 || height != 2000 {
		t.Errorf("Expected dim (%d X %d) got (%d X %d)", 4000, 2000, width, height)
	}
}

func Test_readPixelsRLE(t *testing.T) {

	tests := []struct {
		name   string
		w, h   int
		data   []byte
		err    error
		outLen int
	}{
		{name: "a", w: 10, h: 10, outLen: 300, data: makeSimple(10, 10, []byte{1, 2, 3, 4})},
		{name: "b", w: 10, h: 10, outLen: 1, data: makeSimple(10, 10, []byte{1, 1, 1}), err: newDecoderErr(memoryError, "requires 300 floats but only got 1 floats available")},
		{name: "c", w: 10, h: 10, outLen: 300, data: makeSimple(10, 10, []byte{1, 1, 1}), err: newDecoderErr(readError, "read pixels EOF")},
		{name: "d", w: 10, h: 10, outLen: 300, data: []byte{2, 2, 0x0f, 0xa0}, err: newDecoderErr(formatError, "wrong scanline width, got 4000, but expected 10")},
		{name: "e", w: 100, h: 100, outLen: 30000, data: read(t, "too_short_scanline_data_rle"), err: newDecoderErr(formatError, "not enough data in scanline for rle")},
		{name: "f", w: 100, h: 100, outLen: 30000, data: read(t, "too_short_scanline_data"), err: newDecoderErr(formatError, "not enough data in scanline")},
		{name: "g", w: 100, h: 100, outLen: 30000, data: read(t, "read_pixel_err"), err: newDecoderErr(readError, "read pixels EOF")},
		{name: "h", w: 100, h: 100, outLen: 30000, data: read(t, "41edef25f25b31c67bbe0d8c842e4bdd792e641b-12"), err: newDecoderErr(readError, "readfull #1 EOF")},
		{name: "i", w: 10, h: 10, outLen: 300, data: []byte{2, 2, 0, 10}, err: newDecoderErr(readError, "readfull #2 EOF")},
		{name: "j", w: 100, h: 100, outLen: 30000, data: read(t, "71436d85e41d24e5b70985d22249eaf90d0f392a-6"), err: newDecoderErr(readError, "readfull #3 EOF")},
		{name: "k", w: 7, h: 10, outLen: 210, data: makeSimple(7, 10, []byte{1, 2, 3, 4})},
		{name: "l", w: 8, h: 1, outLen: 24, data: []byte{2, 2, 0, 8, 132, 1, 132, 2, 132, 3, 132, 4, 132, 5, 132, 6, 132, 7, 132, 8}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(subT *testing.T) {
			reader := bytes.NewBuffer(tc.data)
			out := make([]float32, tc.outLen)
			err := readPixelsRLE(reader, tc.w, tc.h, out)

			if tc.err == nil && err != nil {
				subT.Fatal(err)
			}
			if tc.err != nil {
				if err == nil {
					subT.Fatal("Expected error, got none")
				}

				if !errors.Is(err, tc.err) {
					subT.Fatalf("Expected error '%+v', got '%+v'", tc.err, err)
				}
				return
			}

			goldenFile := filepath.Join("testdata", t.Name()+"_"+tc.name+".golden")
			if *update {
				subT.Logf("updating golden file %s", goldenFile)
				gdata := float32ToBytes(subT, out)
				t.Logf("'%s'", gdata.Bytes())
				if err := ioutil.WriteFile(goldenFile, gdata.Bytes(), 0644); err != nil {
					subT.Fatalf("failed to update golden file: %s", err)
				}
			}

			expected, err := ioutil.ReadFile(goldenFile)
			if err != nil {
				subT.Fatalf("failed reading .golden: %s", err)
			}

			actual := float32ToBytes(subT, out)

			if !bytes.Equal(actual.Bytes(), expected) {
				subT.Errorf("readPixelsRLE output does not match .golden file")
			}
		})
	}
}

func makeSimple(w, h int, template []byte) []byte {
	var res []byte
	for i := 0; i < w; i++ {
		for j := 0; j < h; j++ {
			res = append(res, template...)
		}
	}
	return res
}

func read(t *testing.T, name string) []byte {
	d, err := ioutil.ReadFile("testdata/" + name)
	if err != nil {
		t.Fatal(err)
	}
	return d
}

func Test_rgbeToFloat(t *testing.T) {
	expected := [3]float32{0.005859375, 0.009765625, 0.013671875}
	var actual [3]float32
	actual[0], actual[1], actual[2] = rgbeToFloat(1, 2, 3, 128)
	for i, tc := range expected {
		if tc != actual[i] {
			t.Errorf("Expected index %d to be %v, got %v", i, tc, actual[i])
		}
	}
}

func Test_rgbeToFloatZeroExp(t *testing.T) {
	expected := [3]float32{0, 0, 0}
	var actual [3]float32
	actual[0], actual[1], actual[2] = rgbeToFloat(1, 2, 3, 0)
	for i, tc := range expected {
		if tc != actual[i] {
			t.Errorf("Expected index %d to be %v, got %v", i, tc, actual[i])
		}
	}
}

func float32ToBytes(t *testing.T, in []float32) bytes.Buffer {
	var out bytes.Buffer
	if err := binary.Write(&out, binary.BigEndian, in); err != nil {
		t.Fatalf("failed to binary convert []float32 to bytes.Buffer")
	}
	return out
}
