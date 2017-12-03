package technique

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/stojg/graphics/lib/rendering/framebuffer"
	"github.com/stojg/graphics/lib/rendering/primitives"
	"github.com/stojg/graphics/lib/rendering/shader"
)

func BrdfLutTexture() *framebuffer.Texture {

	texture := framebuffer.NewTexture(gl.COLOR_ATTACHMENT0, 512, 512, gl.RG16F, gl.RG, gl.FLOAT, gl.LINEAR, false)

	shad := shader.NewShader("ibl_brdf")
	shad.Bind()

	texture.BindAsRenderTarget()
	texture.SetViewPort()

	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.Disable(gl.DEPTH_TEST)
	primitives.DrawQuad()
	gl.Enable(gl.DEPTH_TEST)
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)

	return texture
}
