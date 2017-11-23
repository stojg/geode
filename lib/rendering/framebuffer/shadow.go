package framebuffer

import "github.com/go-gl/gl/v4.1-core/gl"

func NewShadow(width, height int) *FBO {
	f := &FBO{}
	gl.GenFramebuffers(1, &f.fbo)
	gl.BindFramebuffer(gl.FRAMEBUFFER, f.fbo)
	defer gl.BindFramebuffer(gl.FRAMEBUFFER, 0)

	f.texture = NewTexture(0, gl.DEPTH_COMPONENT16, gl.DEPTH_COMPONENT, gl.FLOAT, width, height)

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)

	borderColor := [4]float32{1.0, 1.0, 1.0, 1.0}
	gl.TexParameterfv(gl.TEXTURE_2D, gl.TEXTURE_BORDER_COLOR, &borderColor[0])

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_BORDER)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_BORDER)

	gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.DEPTH_ATTACHMENT, gl.TEXTURE_2D, f.texture.id, 0)

	gl.DrawBuffer(gl.NONE)
	gl.ReadBuffer(gl.NONE)

	return f
}
