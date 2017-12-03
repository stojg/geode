package framebuffer

import "github.com/go-gl/gl/v4.1-core/gl"

const textType = gl.TEXTURE_2D

func NewTexture(attachment uint32, width int, height int, internalFormat int32, format, xtype uint32, filter int32, clamp bool) *Texture {
	texture := &Texture{
		attachment: attachment,
		width:      int32(width),
		height:     int32(height),
	}
	gl.GenTextures(1, &texture.id)
	gl.BindTexture(textType, texture.id)

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, filter)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)

	if clamp {
		borderColor := [4]float32{1.0, 1.0, 1.0, 1.0}
		gl.TexParameterfv(gl.TEXTURE_2D, gl.TEXTURE_BORDER_COLOR, &borderColor[0])
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_BORDER)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_BORDER)
	} else {
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	}

	if filter == gl.NEAREST_MIPMAP_NEAREST || filter == gl.NEAREST_MIPMAP_LINEAR || filter == gl.LINEAR_MIPMAP_NEAREST || filter == gl.LINEAR_MIPMAP_LINEAR {
		gl.GenerateMipmap(gl.TEXTURE_2D)
	} else {
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_BASE_LEVEL, 0)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAX_LEVEL, 0)
	}

	gl.TexImage2D(gl.TEXTURE_2D, 0, internalFormat, texture.width, texture.height, 0, format, xtype, nil)

	gl.BindTexture(textType, 0)

	// create fbo
	gl.GenFramebuffers(1, &texture.fbo)
	gl.BindFramebuffer(gl.FRAMEBUFFER, texture.fbo)

	gl.FramebufferTexture2D(gl.FRAMEBUFFER, attachment, textType, texture.id, 0)

	if attachment != gl.DEPTH_ATTACHMENT {
		gl.GenRenderbuffers(1, &texture.rbo)
		gl.BindRenderbuffer(gl.RENDERBUFFER, texture.rbo)
		gl.RenderbufferStorage(gl.RENDERBUFFER, gl.DEPTH_COMPONENT16, int32(width), int32(height))
		//gl.RenderbufferStorage(gl.RENDERBUFFER, gl.DEPTH_COMPONENT24, int32(width), int32(height))
		gl.FramebufferRenderbuffer(gl.FRAMEBUFFER, gl.DEPTH_ATTACHMENT, gl.RENDERBUFFER, texture.rbo)
	}

	if gl.CheckFramebufferStatus(gl.FRAMEBUFFER) != gl.FRAMEBUFFER_COMPLETE {
		panic("Framebuffer creation failed, FBO isn't complete.")
	}

	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)

	return texture
}

type Texture struct {
	id         uint32
	attachment uint32
	width      int32
	height     int32

	fbo uint32
	rbo uint32
}

func (t *Texture) Height() int32 {
	return t.height
}

func (t *Texture) Width() int32 {
	return t.width
}

func (t *Texture) ID() uint32 {
	return t.id
}

func (t *Texture) Bind(unit uint32) {
	gl.ActiveTexture(gl.TEXTURE0 + unit)
	gl.BindTexture(textType, t.id)
}

func (t *Texture) BindAsRenderTarget() {
	gl.BindFramebuffer(gl.FRAMEBUFFER, t.fbo)
}

func (t *Texture) SetViewPort() {
	gl.Viewport(0, 0, t.width, t.height)
}
