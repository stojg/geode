package rendering

import (
	"fmt"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/stojg/graphics/lib/debug"
	"github.com/stojg/graphics/lib/loaders"
)

func NewTexture(filename string, srgb bool) *Texture {
	t := &Texture{
		filename: filename,
	}
	resource, err := loadLDRTexture(filename, srgb)
	if err != nil {
		panic(err)
	} else {
		t.resource = resource
	}
	debug.CheckForError("Error during texture load")
	return t
}

func NewMetallicTexture(filename string) *Texture {
	t := &Texture{
		filename: filename,
	}
	resource, err := loadMetallicTexture(filename)
	if err != nil {
		fmt.Println(err)
	} else {
		t.resource = resource
	}
	debug.CheckForError("Error during texture load")
	return t
}

func NewRoughnessTexture(filename string) *Texture {
	t := &Texture{
		filename: filename,
	}
	resource, err := loadRoughnessTexture(filename)
	if err != nil {
		fmt.Println(err)
	} else {
		t.resource = resource
	}
	debug.CheckForError("Error during texture load")
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

func loadLDRTexture(filename string, srgb bool) (*TextureResource, error) {
	rgba, err := loaders.RGBAImagedata(filename)
	if err != nil {
		return nil, err
	}
	rgba = loaders.Flip(rgba)

	if srgb {
		return createTextureResource(rgba.Rect.Size().X, rgba.Rect.Size().Y, gl.SRGB_ALPHA, gl.UNSIGNED_BYTE, rgba.Pix), nil
	} else {
		return createTextureResource(rgba.Rect.Size().X, rgba.Rect.Size().Y, gl.RGBA, gl.UNSIGNED_BYTE, rgba.Pix), nil
	}
}

func loadMetallicTexture(filename string) (*TextureResource, error) {
	rgba, err := loaders.RGBAImagedata(filename)
	if err != nil {
		return nil, err
	}
	rgba = loaders.Flip(rgba)
	return createTextureResource(rgba.Rect.Size().X, rgba.Rect.Size().Y, gl.R8, gl.UNSIGNED_BYTE, rgba.Pix), nil
}

func loadRoughnessTexture(filename string) (*TextureResource, error) {
	rgba, err := loaders.RGBAImagedata(filename)
	if err != nil {
		return nil, err
	}
	rgba = loaders.Flip(rgba)
	return createTextureResource(rgba.Rect.Size().X, rgba.Rect.Size().Y, gl.R8, gl.UNSIGNED_BYTE, rgba.Pix), nil
}

func createTextureResource(width, height int, internalFormat int32, dataType uint32, data []uint8) *TextureResource {
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
	gl.TexImage2D(gl.TEXTURE_2D, 0, internalFormat, resource.width, resource.height, 0, gl.RGBA, dataType, gl.Ptr(data))

	gl.GenerateMipmap(gl.TEXTURE_2D)

	gl.BindTexture(gl.TEXTURE_2D, 0)
	return resource
}
