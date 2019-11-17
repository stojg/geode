package framebuffer

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/stojg/geode/lib/debug"
)

func NewMultiSampledTexture(attachment uint32, width int, height int, internalFormat int32, format, xtype uint32, filter int32, clamp bool) *Texture {
	texture := &Texture{
		width:       int32(width),
		height:      int32(height),
		multiSample: true,
	}
	create(texture, attachment, filter, clamp, internalFormat, format, xtype, width, height)
	return texture
}

func NewTexture(attachment uint32, width int, height int, internalFormat int32, format, xtype uint32, filter int32, clamp bool) *Texture {
	texture := &Texture{
		width:  int32(width),
		height: int32(height),
	}
	create(texture, attachment, filter, clamp, internalFormat, format, xtype, width, height)
	return texture
}

type Texture struct {
	width         int32
	height        int32
	fbo           uint32
	colourTexture uint32
	colourBuffer  uint32
	depthBuffer   uint32
	depthTexture  uint32
	multiSample   bool
}

func (t *Texture) DepthTexture() uint32 {
	return t.depthTexture
}

func (t *Texture) ColourTexture() uint32 {
	return t.colourTexture
}

func (t *Texture) Height() int32 {
	return t.height
}

func (t *Texture) Width() int32 {
	return t.width
}

func (t *Texture) ID() uint32 {
	return t.colourTexture
}

func (t *Texture) Activate(textureSlot uint32) {
	gl.ActiveTexture(gl.TEXTURE0 + textureSlot)
	gl.BindTexture(gl.TEXTURE_2D, t.colourTexture)
}

// Binds the frame buffer, setting it as the current render target. Anything
// rendered after this will be rendered to this FBO, and not to the screen.
func (t *Texture) BindFrameBuffer() {
	gl.BindFramebuffer(gl.FRAMEBUFFER, t.fbo)
	gl.Viewport(0, 0, t.width, t.height)
}

func (t *Texture) UnbindFrameBuffer() {
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
}

func (t *Texture) ResolveToFBO(out *Texture) {
	gl.BindFramebuffer(gl.DRAW_FRAMEBUFFER, out.fbo)
	gl.BindFramebuffer(gl.READ_FRAMEBUFFER, t.fbo)
	gl.BlitFramebuffer(0, 0, t.width, t.height, 0, 0, out.width, out.height, gl.COLOR_BUFFER_BIT|gl.DEPTH_BUFFER_BIT, gl.NEAREST)
	t.UnbindFrameBuffer()
}

func (t *Texture) ResolveToScreen(width, height int32) {
	gl.BindFramebuffer(gl.DRAW_FRAMEBUFFER, 0)
	gl.BindFramebuffer(gl.READ_FRAMEBUFFER, t.fbo)
	gl.DrawBuffer(gl.BACK)
	gl.BlitFramebuffer(0, 0, t.width, t.height, 0, 0, width, height, gl.COLOR_BUFFER_BIT, gl.NEAREST)
	t.UnbindFrameBuffer()
}

func (t *Texture) CleanUp() {
	gl.DeleteFramebuffers(1, &t.fbo)
	gl.DeleteTextures(1, &t.colourTexture)
	//gl.DeleteTextures(1. depthTexture)
	gl.DeleteRenderbuffers(1, &t.depthBuffer)
	gl.DeleteRenderbuffers(1, &t.colourBuffer)
}

func create(texture *Texture, attachment uint32, filter int32, clamp bool, internalFormat int32, format uint32, xtype uint32, width int, height int) {
	gl.GenFramebuffers(1, &texture.fbo)
	gl.BindFramebuffer(gl.FRAMEBUFFER, texture.fbo)

	if texture.multiSample {
		createMultiSampledColourAttachment(texture, width, height)
	} else {
		createTextureAttachment(texture, attachment, filter, clamp, internalFormat, format, xtype)
	}

	if attachment != gl.DEPTH_ATTACHMENT {
		createDepthBufferAttachment(texture, width, height)
	}
	debug.CheckForError("framebuffer.Texture end")
	debug.FramebufferComplete("framebuffer.Texture")

	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
}

// Creates a texture and sets it as the colour buffer attachment for this FBO.
func createTextureAttachment(texture *Texture, attachment uint32, filter int32, clamp bool, internalFormat int32, format uint32, xtype uint32) {
	if texture.width == 0 || texture.height == 0 {
		panic("texture cannot have zero height or width")
	}

	gl.GenTextures(1, &texture.colourTexture)
	gl.BindTexture(gl.TEXTURE_2D, texture.colourTexture)

	gl.TexImage2D(gl.TEXTURE_2D, 0, internalFormat, texture.width, texture.height, 0, format, xtype, nil)
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

	gl.FramebufferTexture2D(gl.FRAMEBUFFER, attachment, gl.TEXTURE_2D, texture.colourTexture, 0)
	gl.BindTexture(gl.TEXTURE_2D, 0)
}

func createMultiSampledColourAttachment(texture *Texture, width int, height int) {
	gl.GenRenderbuffers(1, &texture.colourBuffer)
	gl.BindRenderbuffer(gl.RENDERBUFFER, texture.colourBuffer)
	gl.RenderbufferStorageMultisample(gl.RENDERBUFFER, 4, gl.RGBA16F, int32(width), int32(height))
	gl.FramebufferRenderbuffer(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0, gl.RENDERBUFFER, texture.colourBuffer)
}

// Adds a depth buffer to the FBO in the form of a texture, which can later be sampled.
//func createDepthTextureAttachment(texture *Texture, width, height int) {
//	gl.GenTextures(1, &texture.depthTexture)
//	gl.BindTexture(gl.TEXTURE_2D, texture.depthTexture)
//	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.DEPTH_COMPONENT24, int32(width), int32(height), 0, gl.DEPTH_COMPONENT, gl.FLOAT, nil)
//	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
//	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
//	gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.DEPTH_ATTACHMENT, gl.TEXTURE_2D, texture.depthTexture, 0)
//}

//  * Adds a depth buffer to the FBO in the form of a render buffer. This can't be used for sampling in the shaders.
func createDepthBufferAttachment(texture *Texture, width int, height int) {
	gl.GenRenderbuffers(1, &texture.depthBuffer)
	gl.BindRenderbuffer(gl.RENDERBUFFER, texture.depthBuffer)
	if !texture.multiSample {
		gl.RenderbufferStorage(gl.RENDERBUFFER, gl.DEPTH_COMPONENT, int32(width), int32(height))
	} else {
		gl.RenderbufferStorageMultisample(gl.RENDERBUFFER, 4, gl.DEPTH_COMPONENT, int32(width), int32(height))
	}
	gl.FramebufferRenderbuffer(gl.FRAMEBUFFER, gl.DEPTH_ATTACHMENT, gl.RENDERBUFFER, texture.depthBuffer)
}
