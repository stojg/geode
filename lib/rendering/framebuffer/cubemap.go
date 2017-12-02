package framebuffer

import (
	"fmt"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/stojg/graphics/lib/loaders"
)

func NewCubeMap(attachment uint32) *CubeMap {
	texture := &CubeMap{
		attachment: attachment,
	}

	gl.GenTextures(1, &texture.id)
	gl.BindTexture(gl.TEXTURE_CUBE_MAP, texture.id)

	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_MIN_FILTER, gl.LINEAR)

	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_R, gl.CLAMP_TO_EDGE)

	gl.BindTexture(gl.TEXTURE_CUBE_MAP, 0)

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

func (t *CubeMap) LoadFromFiles(hdr bool, files [6]string) {
	gl.BindTexture(gl.TEXTURE_CUBE_MAP, t.id)
	var internalFormat int32 = gl.SRGB
	if hdr {
		internalFormat = gl.RGBA
	}

	for i := 0; i < 6; i++ {
		img, err := loaders.RGBAImagedata(files[i])
		if err != nil {
			fmt.Println(err)
			continue
		}
		gl.TexImage2D(gl.TEXTURE_CUBE_MAP_POSITIVE_X+uint32(i), 0, internalFormat, int32(img.Rect.Size().X), int32(img.Rect.Size().Y), 0, gl.RGBA, gl.UNSIGNED_INT_8_8_8_8_REV, gl.Ptr(img.Pix))
	}

	gl.BindTexture(gl.TEXTURE_CUBE_MAP, 0)

	gl.GenFramebuffers(1, &t.fbo)
	gl.BindFramebuffer(gl.FRAMEBUFFER, t.fbo)
	gl.FramebufferTexture2D(gl.FRAMEBUFFER, t.attachment, gl.TEXTURE_CUBE_MAP_POSITIVE_X, t.id, 0)

	checkForError("framebuffer.Cubemap end")
	if gl.CheckFramebufferStatus(gl.FRAMEBUFFER) != gl.FRAMEBUFFER_COMPLETE {
		panic("Shadow Framebuffer creation failed, FBO isn't complete.")
	}

	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)

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
