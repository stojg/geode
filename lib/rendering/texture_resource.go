package rendering

import "github.com/go-gl/gl/v4.1-core/gl"

func NewTextureResource() *TextureResource {
	t := &TextureResource{
		refCount: 1,
	}
	gl.GenTextures(1, &t.id)
	return t
}

type TextureResource struct {
	id       uint32
	refCount int
}

func (t *TextureResource) ID() uint32 {
	return t.id
}
