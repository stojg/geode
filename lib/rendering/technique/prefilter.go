package technique

import (
	"math"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/stojg/graphics/lib/components"
	"github.com/stojg/graphics/lib/debug"
	"github.com/stojg/graphics/lib/rendering/primitives"
	"github.com/stojg/graphics/lib/rendering/shader"
)

func Prefilter(src, dest components.Texture) {

	captureProjection := mgl32.Perspective(float32((90*math.Pi)/180.0), 1, 0.1, 10)
	captureViews := []mgl32.Mat4{
		mgl32.LookAt(0, 0, 0, 1, 0, 0, 0, -1, 0),
		mgl32.LookAt(0, 0, 0, -1, 0, 0, 0, -1, 0),
		mgl32.LookAt(0, 0, 0, 0, 1, 0, 0, 0, 1),
		mgl32.LookAt(0, 0, 0, 0, -1, 0, 0, 0, -1),
		mgl32.LookAt(0, 0, 0, 0, 0, 1, 0, -1, 0),
		mgl32.LookAt(0, 0, 0, 0, 0, -1, 0, -1, 0),
	}

	shader := shader.NewShader("ibl_prefilter")
	shader.Bind()

	shader.UpdateUniform("projection", captureProjection)
	shader.UpdateUniform("environmentMap", int32(0))
	src.Activate(0)

	maxMipLevels := 5
	gl.Disable(gl.CULL_FACE)
	for mip := 0; mip < maxMipLevels; mip++ {
		// reisze framebuffer according to mip-level size.
		mipWidth := int32(128 * math.Pow(0.5, float64(mip)))
		mipHeight := int32(128 * math.Pow(0.5, float64(mip)))

		dest.BindFrameBuffer()

		gl.Viewport(0, 0, mipWidth, mipHeight)

		roughness := float32(mip) / float32(maxMipLevels-1)
		shader.UpdateUniform("roughness", roughness)

		for i := 0; i < 6; i++ {
			shader.UpdateUniform("view", captureViews[i])
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
