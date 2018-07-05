package technique

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/stojg/graphics/lib/components"
	"github.com/stojg/graphics/lib/rendering/framebuffer"
	"github.com/stojg/graphics/lib/rendering/primitives"
	"github.com/stojg/graphics/lib/rendering/shader"
)

func NewSkyBox(filename string, s components.RenderState) *SkyBox {
	b := &SkyBox{
		RenderState: s,
		shader:      shader.NewShader("skybox"),
		cubeMap:     framebuffer.NewHDRCubeMap(1024, 1024, filename),
	}

	b.irradianceMap = framebuffer.NewCubeMap(32, 32, false)
	Convolute(b.cubeMap, b.irradianceMap)

	b.preFilterMap = framebuffer.NewCubeMap(128, 128, true)
	Prefilter(b.cubeMap, b.preFilterMap)

	s.AddSamplerSlot("x_skybox")
	s.AddSamplerSlot("x_irradianceMap")
	s.AddSamplerSlot("x_prefilterMap")
	s.AddSamplerSlot("x_brdfLUT")

	b.brdfLut = BrdfLutTexture()

	return b
}

type SkyBox struct {
	components.RenderState
	shader        components.Shader
	cubeMap       *framebuffer.CubeMap
	irradianceMap *framebuffer.CubeMap
	preFilterMap  *framebuffer.CubeMap
	brdfLut       *framebuffer.Texture
}

func (s *SkyBox) CubeMap() *framebuffer.CubeMap {
	return s.cubeMap
}

func (s *SkyBox) IrradianceMap() *framebuffer.CubeMap {
	return s.irradianceMap
}

func (s *SkyBox) PreFilterMap() *framebuffer.CubeMap {
	return s.preFilterMap
}

func (s *SkyBox) BrdfLutTexture() *framebuffer.Texture {
	return s.brdfLut
}

func (s *SkyBox) Load() {
	s.SetTexture("x_irradianceMap", s.irradianceMap)
	s.SetTexture("x_prefilterMap", s.preFilterMap)
	s.SetTexture("x_brdfLUT", s.brdfLut)
}

func (s *SkyBox) Render() {
	gl.DepthFunc(gl.LEQUAL)
	defer gl.DepthFunc(gl.LESS)

	gl.CullFace(gl.FRONT)
	defer gl.CullFace(gl.BACK)

	s.shader.Bind()
	s.RenderState.SetTexture("x_skybox", s.cubeMap)
	s.shader.UpdateUniforms(nil, s.RenderState)
	primitives.DrawCube()
}
