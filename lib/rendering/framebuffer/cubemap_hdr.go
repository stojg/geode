package framebuffer

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/stojg/geode/lib/debug"
)

// NewHDRCubeMap returns a new CubeMap that has been loaded from aHDR Equirectangular image
func NewHDRCubeMap(width, height int32, filename string) *CubeMap {
	texture := &CubeMap{
		attachment: gl.COLOR_ATTACHMENT0,
		width:      width,
		height:     height,
	}
	texture.init(false)
	texture.loadEquiRectangular(filename)

	debug.CheckForError("framebuffer.HDRCubeMap end")
	debug.FramebufferComplete("framebuffer.HDRCubeMap")

	gl.BindTexture(gl.TEXTURE_CUBE_MAP, 0)
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)

	return texture
}
