package rendering

import (
	"fmt"
	"image"
	"image/draw"
	_ "image/png"
	"os"

	"github.com/go-gl/gl/v4.1-core/gl"
)

func NewTexture(filename string) *Texture {

	t := &Texture{
		filename: filename,
	}

	resource, err := LoadTexture(filename)
	if err != nil {
		fmt.Println(err)
	} else {
		t.resource = resource
	}

	return t
}

type Texture struct {
	filename string
	resource *TextureResource
}

func (t *Texture) Bind(slot uint32) {
	gl.ActiveTexture(gl.TEXTURE0 + slot)
	gl.BindTexture(gl.TEXTURE_2D, t.resource.ID())
}

func LoadTexture(filename string) (*TextureResource, error) {
	imgFile, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("texture %q not found on disk: %v", filename, err)
	}
	img, _, err := image.Decode(imgFile)
	if err != nil {
		return nil, err
	}

	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		return nil, fmt.Errorf("unsupported stride %d", rgba.Stride)
	}
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)

	// flip it into open GL format
	rgba = flip(rgba)

	//var internalFormat int32 = gl.RGBA
	//if gammaCorrect {
	//	internalFormat = gl.SRGB8_ALPHA8
	//}

	resource := NewTextureResource()
	gl.BindTexture(gl.TEXTURE_2D, resource.ID())
	gl.ActiveTexture(gl.TEXTURE0)

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)

	gl.TexParameterf(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameterf(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)

	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, int32(rgba.Rect.Size().X), int32(rgba.Rect.Size().Y), 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(rgba.Pix))

	return resource, nil
}

// flip the image upside down so that OpenGL can use it as a texture properly
func flip(src *image.RGBA) *image.RGBA {
	srcW := src.Bounds().Max.X
	srcH := src.Bounds().Max.Y
	dstW := srcW
	dstH := srcH

	dst := image.NewRGBA(src.Bounds())

	for dstY := 0; dstY < dstH; dstY++ {
		for dstX := 0; dstX < dstW; dstX++ {
			srcX := dstX
			srcY := dstH - dstY - 1
			srcOff := srcY*src.Stride + srcX*4
			dstOff := dstY*dst.Stride + dstX*4
			copy(dst.Pix[dstOff:dstOff+4], src.Pix[srcOff:srcOff+4])
		}
	}

	return dst
}
