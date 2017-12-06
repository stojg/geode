package framebuffer

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/stojg/graphics/lib/debug"
)

func NewLDRCubeMap(files [6]string) *CubeMap {
	texture := &CubeMap{
		attachment: gl.COLOR_ATTACHMENT0,
	}
	initCubeMap(texture, false)
	loadFromFiles(texture, files)

	debug.CheckForError("framebuffer.LDRCubeMap end")
	debug.FramebufferComplete("framebuffer.LDRCubeMap")

	gl.BindTexture(gl.TEXTURE_CUBE_MAP, 0)
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)

	return texture
}
