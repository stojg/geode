package technique

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/stojg/graphics/lib/components"
	"github.com/stojg/graphics/lib/rendering/framebuffer"
	"github.com/stojg/graphics/lib/rendering/primitives"
	"github.com/stojg/graphics/lib/rendering/shader"
)

// Convolute takes a cubemap texture and convolutes
func Convolute(src, dest components.Texture) {

	shad := shader.NewShader("convolute")
	shad.Bind()

	shad.UpdateUniform("projection", framebuffer.CubeProjection())
	shad.UpdateUniform("environmentMap", int32(0))
	src.Activate(0)

	dest.BindFrameBuffer()

	gl.Disable(gl.CULL_FACE)
	for i, captureView := range framebuffer.CubeViews() {
		shad.UpdateUniform("view", captureView)
		gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0, gl.TEXTURE_CUBE_MAP_POSITIVE_X+uint32(i), dest.ID(), 0)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		primitives.DrawCube()
	}
	gl.Enable(gl.CULL_FACE)
}
