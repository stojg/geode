package loaders

import (
	"fmt"
	"image"
	"image/draw"
	_ "image/jpeg"
	_ "image/png"
	"os"
)

func RGBAImagedata(filename string) (*image.RGBA, error) {
	imgFile, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("texture %q not found on disk: %v", filename, err)
	}
	defer imgFile.Close()
	img, _, err := image.Decode(imgFile)
	if err != nil {
		return nil, err
	}

	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		return nil, fmt.Errorf("unsupported stride %d", rgba.Stride)
	}
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)

	return rgba, nil
}

func RGBEImagedata(filename string) (int, int, []float32, error) {
	imgFile, err := os.Open(filename)
	if err != nil {
		return 0, 0, nil, fmt.Errorf("texture %q not found on disk: %v", filename, err)
	}
	defer imgFile.Close()
	return DecodeRGBE(imgFile)
}

// Flip the image upside down so that OpenGL can use it as a texture properly
func Flip(src *image.RGBA) *image.RGBA {
	maxX := src.Bounds().Max.X
	maxY := src.Bounds().Max.Y

	dst := image.NewRGBA(src.Bounds())

	for y := 0; y < maxY; y++ {
		for x := 0; x < maxX; x++ {
			srcY := maxY - y - 1
			srcRow := srcY*src.Stride + x*4
			destRow := y*src.Stride + x*4
			copy(dst.Pix[destRow:destRow+4], src.Pix[srcRow:srcRow+4])
		}
	}

	return dst
}

func FlipRaw(width, height int, src []float32) []float32 {
	dst := make([]float32, len(src))

	rowSize := width * 3

	for y := 0; y < height; y++ {
		srcStart := y * rowSize
		srcEnd := srcStart + rowSize

		dstStart := (height - y - 1) * rowSize
		dstEnd := dstStart + rowSize

		copy(dst[dstStart:dstEnd], src[srcStart:srcEnd])
	}
	return dst
}
