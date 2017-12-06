package framebuffer

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/stojg/graphics/lib/debug"
)

func NewCubeMap(width, height int32, mipMap bool) *CubeMap {
	texture := &CubeMap{
		attachment: gl.COLOR_ATTACHMENT0,
		width:      width,
		height:     height,
	}
	initCubeMap(texture, mipMap)
	reserveCubeMap(texture)

	if mipMap {
		gl.GenerateMipmap(gl.TEXTURE_CUBE_MAP)
	}

	debug.CheckForError("framebuffer.CubeMap end")
	debug.FramebufferComplete("framebuffer.LDRCubeMap")

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
	gl.BindFramebuffer(gl.FRAMEBUFFER, t.fbo)
}

func (t *CubeMap) SetViewPort() {
	gl.Viewport(0, 0, t.width, t.height)
}

func initCubeMap(t *CubeMap, mipMap bool) {
	gl.GenFramebuffers(1, &t.fbo)
	gl.GenRenderbuffers(1, &t.rbo)

	gl.BindFramebuffer(gl.FRAMEBUFFER, t.fbo)
	gl.BindRenderbuffer(gl.RENDERBUFFER, t.rbo)

	gl.RenderbufferStorage(gl.RENDERBUFFER, gl.DEPTH_COMPONENT24, t.width, t.height)
	gl.FramebufferRenderbuffer(gl.FRAMEBUFFER, gl.DEPTH_ATTACHMENT, gl.RENDERBUFFER, t.rbo)

	gl.GenTextures(1, &t.id)

	gl.BindTexture(gl.TEXTURE_CUBE_MAP, t.id)
	gl.FramebufferTexture2D(gl.FRAMEBUFFER, t.attachment, gl.TEXTURE_CUBE_MAP_POSITIVE_X, t.id, 0)

	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_R, gl.CLAMP_TO_EDGE)

	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	if mipMap {
		gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_MIN_FILTER, gl.LINEAR_MIPMAP_LINEAR)
	} else {
		gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
		gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_BASE_LEVEL, 0)
		gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_MAX_LEVEL, 0)
	}

	if debug.CheckForError("initCubeMap") {
		panic("init cubemap")
	}
}
