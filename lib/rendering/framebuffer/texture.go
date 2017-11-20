package framebuffer

import "github.com/go-gl/gl/v4.1-core/gl"

func NewTexture(attachment uint32, format int32, xtype uint32, width int32, height int32) *Texture {
	texture := &Texture{
		attachment: gl.COLOR_ATTACHMENT0 + attachment,
	}
	gl.GenTextures(1, &texture.id)
	gl.BindTexture(gl.TEXTURE_2D, texture.id)
	gl.TexImage2D(gl.TEXTURE_2D, 0, format, width, height, 0, gl.RGBA, xtype, nil)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.FramebufferTexture2D(gl.FRAMEBUFFER, texture.attachment, gl.TEXTURE_2D, texture.id, 0)
	return texture
}

type Texture struct {
	id         uint32
	attachment uint32
}

func (t *Texture) Bind() {
	gl.BindTexture(gl.TEXTURE_2D, t.id)
}

func (t *Texture) DrawInto() {
	gl.DrawBuffer(t.attachment)
}
