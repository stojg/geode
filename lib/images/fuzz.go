package images

import "bytes"

func Fuzz(data []byte) int {
	reader := bytes.NewBuffer(data)
	out := make([]float32, 10*10*3)
	err := readPixelsRLE(reader, 100, 100, out)
	if err != nil {
		return -1
	}
	return 1
}
