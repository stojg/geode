package lights

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/stojg/graphics/lib/components"
	"github.com/stojg/graphics/lib/rendering"
	"github.com/stojg/graphics/lib/rendering/framebuffer"
)

func NewDirectional(r, g, b, intensity float32) *Directional {
	return &Directional{
		BaseLight: BaseLight{
			color:  mgl32.Vec3{r, g, b}.Mul(intensity),
			shader: rendering.NewShader("forward_directional"),
		},
	}
}

type Directional struct {
	BaseLight

	samplerName string
	textureSlot uint32
	texture     *framebuffer.Texture
}

func (b *Directional) AddToEngine(e components.Engine) {
	e.GetRenderingEngine().AddLight(b)
}

func (b *Directional) Direction() mgl32.Vec3 {
	return b.BaseLight.Position()
}

func (b *Directional) SetShadowTexture(slot uint32, samplerName string, texture *framebuffer.Texture) {
	b.samplerName = samplerName
	b.textureSlot = slot
	b.texture = texture
}

func (b *Directional) BindShadow() {
	gl.ActiveTexture(gl.TEXTURE0 + uint32(b.textureSlot))
	b.Shader().SetUniform(b.samplerName, int32(b.textureSlot))
	gl.BindTexture(gl.TEXTURE_2D, b.texture.ID())
}
