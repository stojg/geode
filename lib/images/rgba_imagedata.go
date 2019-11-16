package images

import (
	"fmt"
	"image"
	"image/draw"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"
)

func RGBAImagedata(filename string) (*image.RGBA, error) {
	imgFile, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("image %q not found on disk: %v", filename, err)
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
func Flip(src *image.RGBA) {
	height := src.Bounds().Dy()

	stride := src.Stride
	scratchBuffer := make([]uint8, stride)
	for y := 0; y < height/2; y++ {
		top := y * stride
		bottom := (height - y - 1) * stride
		// copy bottom row to buffer
		copy(scratchBuffer[:], src.Pix[bottom:bottom+stride])
		// copy top row to bottom row
		copy(src.Pix[bottom:bottom+stride], src.Pix[top:top+stride])
		// copy buffer (previous bottom) to top row
		copy(src.Pix[top:top+stride], scratchBuffer[:])
	}
}

func FlipRaw(width, height int, src []float32) {
	const valuesPerPixel = 3

	if width*height*valuesPerPixel != len(src) {
		log.Fatalf("width * height (%d) doesn't add up length of src (%d)", width*height*valuesPerPixel, len(src))
	}

	stride := width * valuesPerPixel
	scratchBuffer := make([]float32, stride)
	for y := 0; y < height/2; y++ {
		top := y * stride
		bottom := (height - y - 1) * stride
		// copy bottom row to buffer
		copy(scratchBuffer[:], src[bottom:bottom+stride])
		// copy top row to bottom row
		copy(src[bottom:bottom+stride], src[top:top+stride])
		// copy buffer (bottom) to top row
		copy(src[top:top+stride], scratchBuffer[:])
	}
}
