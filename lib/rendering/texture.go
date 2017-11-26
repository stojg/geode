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

func (t *Texture) Height() int32 {
	return t.resource.height
}

func (t *Texture) Width() int32 {
	return t.resource.width
}

func (t *Texture) Bind(slot uint32) {
	gl.ActiveTexture(gl.TEXTURE0 + slot)
	gl.BindTexture(gl.TEXTURE_2D, t.resource.ID())
}

func (t *Texture) SetViewPort() {
	gl.Viewport(0, 0, t.Width(), t.Height())
}
func (t *Texture) BindAsRenderTarget() {
	panic("Cant write to material textures you mad lad!")
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

	// @todo this is handy when dealing with HDR/nono SRGB image sources
	//var internalFormat int32 = gl.RGBA
	//if gammaCorrect {
	//	internalFormat = gl.SRGB8_ALPHA8
	//}

	resource := NewTextureResource()
	resource.width = int32(rgba.Rect.Size().X)
	resource.height = int32(rgba.Rect.Size().Y)
	gl.BindTexture(gl.TEXTURE_2D, resource.ID())
	gl.ActiveTexture(gl.TEXTURE0)

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)

	gl.TexParameterf(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR_MIPMAP_LINEAR)
	gl.TexParameterf(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)

	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, resource.width, resource.height, 0, gl.RGBA, gl.UNSIGNED_INT_8_8_8_8_REV, gl.Ptr(rgba.Pix))

	gl.GenerateMipmap(gl.TEXTURE_2D)

	return resource, nil
}

// flip the image upside down so that OpenGL can use it as a texture properly
func flip(src *image.RGBA) *image.RGBA {
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
