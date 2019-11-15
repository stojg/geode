package technique

import (
	"math"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/stojg/geode/lib/components"
	"github.com/stojg/geode/lib/debug"
	"github.com/stojg/geode/lib/rendering/framebuffer"
	"github.com/stojg/geode/lib/rendering/primitives"
	"github.com/stojg/geode/lib/rendering/shader"
)

func Prefilter(src, dest components.Texture) {

	shader := shader.NewShader("ibl_prefilter")
	shader.Bind()

	shader.UpdateUniform("projection", framebuffer.CubeProjection())
	shader.UpdateUniform("environmentMap", int32(0))
	src.Activate(0)

	maxMipLevels := 5
	gl.Disable(gl.CULL_FACE)
	for mip := 0; mip < maxMipLevels; mip++ {
		// reisze framebuffer according to mip-level size.
		mipWidth := int32(float64(dest.Width()) * math.Pow(0.5, float64(mip)))
		mipHeight := int32(float64(dest.Height()) * math.Pow(0.5, float64(mip)))

		dest.BindFrameBuffer()

		gl.Viewport(0, 0, mipWidth, mipHeight)

		roughness := float32(mip) / float32(maxMipLevels-1)
		shader.UpdateUniform("roughness", roughness)

		for i, captureView := range framebuffer.CubeViews() {
			shader.UpdateUniform("view", captureView)
			gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0, gl.TEXTURE_CUBE_MAP_POSITIVE_X+uint32(i), dest.ID(), int32(mip))
			gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
			primitives.DrawCube()
		}

	}
	gl.Enable(gl.CULL_FACE)

	if debug.CheckForError("prefilter - end") {
		panic("crap")
	}

}
