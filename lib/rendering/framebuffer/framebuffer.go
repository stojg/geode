package framebuffer

import "github.com/go-gl/gl/v4.1-core/gl"

type FBO struct {
	fbo     uint32
	rbo     uint32
	texture *Texture
	vao     uint32
}

func (f *FBO) Bind() {
	gl.BindFramebuffer(gl.FRAMEBUFFER, f.fbo)
	f.texture.DrawInto()
}

func (f *FBO) BindTexture() {
	f.texture.Bind()
}

func (f *FBO) Texture() *Texture {
	return f.texture
}
