package images

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"flag"
	"fmt"
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
	t.Log(r)
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
		name string
		w, h int
		data []byte
		err  error
	}{
		{name: "basic", w: 10, h: 10, data: []byte{1, 2, 3, 4}},
		{name: "corrupt", w: 10, h: 10, data: []byte{99, 99, 99}, err: newError(memoryError, errors.New(""))},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(subT *testing.T) {
			var in []byte
			for i := 0; i < tc.w; i++ {
				for j := 0; j < tc.h; j++ {
					in = append(in, tc.data...)
				}
			}

			reader := bytes.NewBuffer(in)
			out := make([]float32, tc.w*tc.h*3)
			err := readPixelsRLE(reader, tc.w, tc.h, out)

			if tc.err == nil && err != nil {
				subT.Fatal(err)
			}
			if tc.err != nil {
				if err == nil {
					subT.Fatal("Expected error, got none")
				}

				if !errors.Is(err, tc.err) {
					subT.Fatalf("Expected error %+v, got %+v", tc.err, err)
				}
				return
			}

			goldenFile := filepath.Join("testdata", t.Name()+"_"+tc.name+".golden")
			if *update {
				subT.Logf("updating golden file %s", goldenFile)
				gdata := float32ToBytes(subT, out)

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

func Test_rgbeToFloat(t *testing.T) {
	expected := [3]float32{0.005859375, 0.009765625, 0.013671875}
	var actual [3]float32
	actual[0], actual[1], actual[2] = rgbeToFloat(1, 2, 3, 128)
	fmt.Println(actual)
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
