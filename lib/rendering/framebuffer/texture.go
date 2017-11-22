package framebuffer

import "github.com/go-gl/gl/v4.1-core/gl"

func NewTexture(attachment uint32, internalformat int32, format, xtype uint32, width int, height int) *Texture {
	texture := &Texture{
		attachment: gl.COLOR_ATTACHMENT0 + attachment,
		width:      int32(width),
		height:     int32(height),
	}
	gl.GenTextures(1, &texture.id)
	gl.BindTexture(gl.TEXTURE_2D, texture.id)
	gl.TexImage2D(gl.TEXTURE_2D, 0, internalformat, texture.width, texture.height, 0, format, xtype, nil)
	return texture
}

type Texture struct {
	id         uint32
	attachment uint32
	width      int32
	height     int32
}

func (t *Texture) ID() uint32 {
	return t.id
}

func (t *Texture) Bind() {
	gl.BindTexture(gl.TEXTURE_2D, t.id)
}

func (t *Texture) DrawInto() {
	gl.DrawBuffer(t.attachment)
}

func (t *Texture) SetViewPort() {
	gl.Viewport(0, 0, t.width, t.height)
}
