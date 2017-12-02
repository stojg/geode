package framebuffer

import (
	"fmt"

	"github.com/go-gl/gl/v4.1-core/gl"
)

func NewCubeMap(width, height int32) *CubeMap {
	texture := &CubeMap{
		attachment: gl.COLOR_ATTACHMENT0,
		width:      width,
		height:     height,
	}
	texture.initCubemap()

	gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0, gl.TEXTURE_CUBE_MAP_POSITIVE_X, texture.id, 0)

	for i := 0; i < 6; i++ {
		gl.TexImage2D(gl.TEXTURE_CUBE_MAP_POSITIVE_X+uint32(i), 0, gl.RGB16F, texture.width, texture.height, 0, gl.RGB, gl.FLOAT, nil)
		checkForError("framebuffer.Cubemap end")
	}

	if e := gl.CheckFramebufferStatus(gl.FRAMEBUFFER); e != gl.FRAMEBUFFER_COMPLETE {
		panic(fmt.Sprintf("Cubmap  Framebuffer creation failed, FBO isn't complete: 0x%x", e))
	}

	gl.BindTexture(gl.TEXTURE_CUBE_MAP, 0)
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)

	return texture
}

type CubeMap struct {
	id         uint32
	attachment uint32
	width      int32
	height     int32

	fbo uint32
	rbo uint32
}

func (t *CubeMap) Height() int32 {
	return t.height
}

func (t *CubeMap) Width() int32 {
	return t.width
}

func (t *CubeMap) ID() uint32 {
	return t.id
}

func (t *CubeMap) Bind(unit uint32) {
	gl.ActiveTexture(gl.TEXTURE0 + unit)
	gl.BindTexture(gl.TEXTURE_CUBE_MAP, t.id)
}

func (t *CubeMap) BindAsRenderTarget() {
	// this probably wont work without setting which side of the cubemap to render to?
	gl.BindTexture(textType, t.id)
	gl.BindFramebuffer(gl.FRAMEBUFFER, t.fbo)
}

func (t *CubeMap) SetViewPort() {
	gl.Viewport(0, 0, t.width, t.height)
}

func (t *CubeMap) initCubemap() {
	gl.GenFramebuffers(1, &t.fbo)
	gl.GenRenderbuffers(1, &t.rbo)

	gl.BindFramebuffer(gl.FRAMEBUFFER, t.fbo)
	gl.BindRenderbuffer(gl.RENDERBUFFER, t.rbo)

	gl.RenderbufferStorage(gl.RENDERBUFFER, gl.DEPTH_COMPONENT24, t.width, t.height)
	gl.FramebufferRenderbuffer(gl.FRAMEBUFFER, gl.DEPTH_ATTACHMENT, gl.RENDERBUFFER, t.rbo)

	gl.GenTextures(1, &t.id)

	gl.BindTexture(gl.TEXTURE_CUBE_MAP, t.id)
	gl.FramebufferTexture2D(gl.FRAMEBUFFER, t.attachment, gl.TEXTURE_CUBE_MAP_POSITIVE_X, t.id, 0)

	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_BASE_LEVEL, 0)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_MAX_LEVEL, 0)

	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_R, gl.CLAMP_TO_EDGE)

	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
}
