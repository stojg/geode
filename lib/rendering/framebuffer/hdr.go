package framebuffer

import (
	"fmt"

	"github.com/go-gl/gl/v4.1-core/gl"
)

func NewHDR(width, height int) *FBO {
	f := &FBO{}
	gl.GenFramebuffers(1, &f.fbo)
	gl.BindFramebuffer(gl.FRAMEBUFFER, f.fbo)
	defer gl.BindFramebuffer(gl.FRAMEBUFFER, 0)

	f.texture = NewTexture(0, gl.RGBA16F, gl.RGBA, gl.FLOAT, width, height)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)

	gl.FramebufferTexture2D(gl.FRAMEBUFFER, f.texture.attachment, gl.TEXTURE_2D, f.texture.id, 0)

	gl.GenRenderbuffers(1, &f.rbo)
	gl.BindRenderbuffer(gl.RENDERBUFFER, f.rbo)
	defer gl.BindRenderbuffer(gl.RENDERBUFFER, 0)
	gl.RenderbufferStorage(gl.RENDERBUFFER, gl.DEPTH24_STENCIL8, int32(width), int32(height))
	gl.FramebufferRenderbuffer(gl.FRAMEBUFFER, gl.DEPTH_STENCIL_ATTACHMENT, gl.RENDERBUFFER, f.rbo)

	if gl.CheckFramebufferStatus(gl.FRAMEBUFFER) != gl.FRAMEBUFFER_COMPLETE {
		fmt.Println("FBO not complete")
	}

	return f
}
