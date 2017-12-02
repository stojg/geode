package technique

import (
	"math"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/stojg/graphics/lib/components"
	"github.com/stojg/graphics/lib/rendering/primitives"
	"github.com/stojg/graphics/lib/rendering/shader"
)

func Convolute(src, dest components.Texture) {

	captureProjection := mgl32.Perspective(float32((90*math.Pi)/180.0), 1, 0.1, 10)
	captureViews := []mgl32.Mat4{
		mgl32.LookAt(0, 0, 0, 1, 0, 0, 0, -1, 0),
		mgl32.LookAt(0, 0, 0, -1, 0, 0, 0, -1, 0),
		mgl32.LookAt(0, 0, 0, 0, 1, 0, 0, 0, 1),
		mgl32.LookAt(0, 0, 0, 0, -1, 0, 0, 0, -1),
		mgl32.LookAt(0, 0, 0, 0, 0, 1, 0, -1, 0),
		mgl32.LookAt(0, 0, 0, 0, 0, -1, 0, -1, 0),
	}

	shad := shader.NewShader("convolute")
	shad.Bind()

	shad.UpdateUniform("projection", captureProjection)
	shad.UpdateUniform("environmentMap", int32(0))
	src.Bind(0)

	dest.SetViewPort()
	dest.BindAsRenderTarget()

	gl.Disable(gl.CULL_FACE)
	for i := 0; i < 6; i++ {
		shad.UpdateUniform("view", captureViews[i])
		gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0, gl.TEXTURE_CUBE_MAP_POSITIVE_X+uint32(i), dest.ID(), 0)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		primitives.DrawCube()
	}
	gl.Enable(gl.CULL_FACE)
}
