package framebuffer

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/stojg/graphics/lib/loaders"
)

func loadHDRTextureResource(filename string) (*TextureResource, error) {
	width, height, data, err := loaders.RGBEImagedata(filename)
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

	gl.TexParameterf(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameterf(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)

	if resource.width == 0 || resource.height == 0 {
		panic("texture cannot have zero height or width")
	}
	gl.TexImage2D(gl.TEXTURE_2D, 0, internalFormat, resource.width, resource.height, 0, gl.RGB, dataType, gl.Ptr(data))

	gl.BindTexture(textType, 0)

	return resource
}

func NewTextureResource() *TextureResource {
	t := &TextureResource{
		refCount: 1,
	}
	gl.GenTextures(1, &t.id)
	return t
}

type TextureResource struct {
	id            uint32
	refCount      int
	width, height int32
}

func (t *TextureResource) ID() uint32 {
	return t.id
}

func (t *TextureResource) Delete() {
	gl.DeleteTextures(1, &t.id)
}
