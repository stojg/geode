package framebuffer

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/stojg/graphics/lib/debug"
)

// NewLDRCubeMap returns a new CubeMap that has been loaded from six LDR images
func NewLDRCubeMap(files [6]string) *CubeMap {
	texture := &CubeMap{
		attachment: gl.COLOR_ATTACHMENT0,
	}
	texture.init(false)
	texture.loadFromFiles(files)

	debug.CheckForError("framebuffer.LDRCubeMap end")
	debug.FramebufferComplete("framebuffer.LDRCubeMap")

	gl.BindTexture(gl.TEXTURE_CUBE_MAP, 0)
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)

	return texture
}
