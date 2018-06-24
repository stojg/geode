package technique

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/stojg/graphics/lib/components"
	"github.com/stojg/graphics/lib/rendering/primitives"
	"github.com/stojg/graphics/lib/rendering/shader"
)

func NewSkyBox(cubemap components.Texture) *SkyBox {
	return &SkyBox{
		shader:         shader.NewShader("skybox"),
		cubemapTexture: cubemap,
	}
}

type SkyBox struct {
	cubemapTexture components.Texture
	shader         components.Shader
}

func (e *SkyBox) Draw(r components.RenderingEngine) {
	gl.DepthFunc(gl.LEQUAL)
	defer gl.DepthFunc(gl.LESS)

	gl.CullFace(gl.FRONT)
	defer gl.CullFace(gl.BACK)

	e.shader.Bind()
	r.SetTexture("x_skybox", e.cubemapTexture)
	r.SetSamplerSlot("x_skybox", 0)
	e.shader.UpdateUniforms(nil, r)
	primitives.DrawCube()
}
