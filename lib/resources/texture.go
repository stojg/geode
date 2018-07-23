package resources

import (
	"fmt"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/stojg/graphics/lib/debug"
	"github.com/stojg/graphics/lib/images"
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

func (t *Texture) Activate(textureSlot uint32) {
	gl.ActiveTexture(gl.TEXTURE0 + textureSlot)
	gl.BindTexture(gl.TEXTURE_2D, t.resource.ID())
}

func (t *Texture) BindFrameBuffer() {
	panic("Cant write to material textures you mad lad!")
}

func (t *Texture) UnbindFrameBuffer() {
	panic("Cant unbind non-FBO texture, you mad lad!")
}

func (t *Texture) SetViewPort() {
	gl.Viewport(0, 0, t.Width(), t.Height())
}

func loadLDRTexture(filename string, srgb bool) (*TextureResource, error) {
	rgba, err := images.RGBAImagedata(filename)
	if err != nil {
		return nil, err
	}
	images.Flip(rgba)

	if srgb {
		return createTextureResource(rgba.Rect.Size().X, rgba.Rect.Size().Y, gl.SRGB_ALPHA, gl.UNSIGNED_BYTE, rgba.Pix), nil
	} else {
		return createTextureResource(rgba.Rect.Size().X, rgba.Rect.Size().Y, gl.RGBA8, gl.UNSIGNED_BYTE, rgba.Pix), nil
	}
}

func loadMetallicTexture(filename string) (*TextureResource, error) {
	rgba, err := images.RGBAImagedata(filename)
	if err != nil {
		return nil, err
	}
	images.Flip(rgba)
	return createTextureResource(rgba.Rect.Size().X, rgba.Rect.Size().Y, gl.R8, gl.UNSIGNED_BYTE, rgba.Pix), nil
}

func loadRoughnessTexture(filename string) (*TextureResource, error) {
	rgba, err := images.RGBAImagedata(filename)
	if err != nil {
		return nil, err
	}
	images.Flip(rgba)
	return createTextureResource(rgba.Rect.Size().X, rgba.Rect.Size().Y, gl.R8, gl.UNSIGNED_BYTE, rgba.Pix), nil
}

func createTextureResource(width, height int, internalFormat int32, dataType uint32, data []uint8) *TextureResource {
	if width == 0 || height == 0 {
		panic("texture cannot have zero height or width")
	}
	resource := NewTextureResource()
	resource.width = int32(width)
	resource.height = int32(height)
	gl.BindTexture(gl.TEXTURE_2D, resource.ID())
	gl.ActiveTexture(gl.TEXTURE0)

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR_MIPMAP_LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)

	if glfw.ExtensionSupported("GL_EXT_texture_filter_anisotropic") {
		var t float32
		gl.GetFloatv(gl.MAX_TEXTURE_MAX_ANISOTROPY, &t)
		if t > 4 {
			t = 4
		}
		gl.TexParameterf(gl.TEXTURE_2D, gl.TEXTURE_MAX_ANISOTROPY, t)
	}

	gl.TexImage2D(gl.TEXTURE_2D, 0, internalFormat, resource.width, resource.height, 0, gl.RGBA, dataType, gl.Ptr(data))

	gl.GenerateMipmap(gl.TEXTURE_2D)

	gl.BindTexture(gl.TEXTURE_2D, 0)
	return resource
}
