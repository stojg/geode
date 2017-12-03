package framebuffer

import (
	"fmt"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/stojg/graphics/lib/loaders"
)

func NewLDRCubeMap(files [6]string) *LDRCubeMap {
	texture := &LDRCubeMap{
		attachment: gl.COLOR_ATTACHMENT0,
		width:      1,
		height:     1,
	}
	texture.initCubemap()
	texture.LoadFromFiles(files)

	return texture
}

type LDRCubeMap struct {
	id         uint32
	attachment uint32
	width      int32
	height     int32

	fbo uint32
	rbo uint32
}

func (t *LDRCubeMap) Height() int32 {
	return t.height
}

func (t *LDRCubeMap) LoadFromFiles(files [6]string) {
	img, err := loaders.RGBAImagedata(files[0])
	if err != nil {
		fmt.Println(err)
		return
	}
	t.width = int32(img.Rect.Size().X)
	t.height = int32(img.Rect.Size().Y)

	for i := 0; i < 6; i++ {
		img, err := loaders.RGBAImagedata(files[i])
		if err != nil {
			fmt.Println(err)
			continue
		}
		gl.TexImage2D(gl.TEXTURE_CUBE_MAP_POSITIVE_X+uint32(i), 0, gl.SRGB, t.width, t.height, 0, gl.RGBA, gl.UNSIGNED_INT_8_8_8_8_REV, gl.Ptr(img.Pix))
	}

	checkForError("framebuffer.Cubemap end")
	if e := gl.CheckFramebufferStatus(gl.FRAMEBUFFER); e != gl.FRAMEBUFFER_COMPLETE {
		panic(fmt.Sprintf("Cubmap LDR Framebuffer creation failed, FBO isn't complete: 0x%x", e))
	}

	gl.BindTexture(gl.TEXTURE_CUBE_MAP, 0)
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
}

func (t *LDRCubeMap) Width() int32 {
	return t.width
}

func (t *LDRCubeMap) ID() uint32 {
	return t.id
}

func (t *LDRCubeMap) Bind(unit uint32) {
	gl.ActiveTexture(gl.TEXTURE0 + unit)
	gl.BindTexture(gl.TEXTURE_CUBE_MAP, t.id)
}

func (t *LDRCubeMap) BindAsRenderTarget() {
	gl.BindFramebuffer(gl.FRAMEBUFFER, t.fbo)
}

func (t *LDRCubeMap) SetViewPort() {
	gl.Viewport(0, 0, t.width, t.height)
}

func (t *LDRCubeMap) initCubemap() {
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
