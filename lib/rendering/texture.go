package rendering

import (
	"fmt"
	"os"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/stojg/graphics/lib/loaders"
)

func NewTexture(filename string) *Texture {
	t := &Texture{
		filename: filename,
	}
	resource, err := loadLDRTexture(filename)
	if err != nil {
		fmt.Println(err)
	} else {
		t.resource = resource
	}
	return t
}

func NewHDRTexture(filename string) *Texture {
	t := &Texture{
		filename: filename,
	}
	resource, err := loadHDRTexture(filename)
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

func (t *Texture) ID() uint32 {
	return t.resource.ID()
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

func loadLDRTexture(filename string) (*TextureResource, error) {
	rgba, err := loaders.RGBAImagedata(filename)
	if err != nil {
		return nil, err
	}
	rgba = loaders.Flip(rgba)
	return createTextureResource(rgba.Rect.Size().X, rgba.Rect.Size().Y, gl.RGB, gl.UNSIGNED_INT_8_8_8_8_REV, rgba.Pix), nil
}

func loadHDRTexture(filename string) (*TextureResource, error) {
	imgFile, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("texture %q not found on disk: %v", filename, err)
	}
	defer imgFile.Close()
	width, height, data, err := loaders.RGBEDecoder(imgFile)
	if err != nil {
		return nil, err
	}
	data = loaders.FlipRaw(width, height, data)
	return createTextureResource(width, height, gl.RGB32F, gl.FLOAT, data), nil
}

func createTextureResource(width, height int, internalFormat int32, dataType uint32, data interface{}) *TextureResource {
	resource := NewTextureResource()
	resource.width = int32(width)
	resource.height = int32(height)
	gl.BindTexture(gl.TEXTURE_2D, resource.ID())
	gl.ActiveTexture(gl.TEXTURE0)

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)

	gl.TexParameterf(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR_MIPMAP_LINEAR)
	gl.TexParameterf(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)

	if width == 0 || height == 0 {
		panic("texture cannot have zero height or width")
	}
	gl.TexImage2D(gl.TEXTURE_2D, 0, internalFormat, resource.width, resource.height, 0, gl.RGB, dataType, gl.Ptr(data))

	gl.GenerateMipmap(gl.TEXTURE_2D)
	return resource
}
