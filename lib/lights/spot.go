package lights

import (
	"math"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/stojg/graphics/lib/components"
	"github.com/stojg/graphics/lib/rendering"
	"github.com/stojg/graphics/lib/rendering/framebuffer"
)

func NewSpot(r, g, b, intensity, angle float32) *Spot {

	radians := float32(math.Cos(float64(mgl32.DegToRad(angle))))

	return &Spot{
		BaseLight: BaseLight{
			color:  mgl32.Vec3{r, g, b}.Mul(intensity),
			shader: rendering.NewShader("forward_spot"),
		},
		PointLight: PointLight{
			constant: 1,
			linear:   0.22,
			exponent: 0.20,
		},
		direction: mgl32.Vec3{0, 0, 0},
		cutoff:    radians,

		shadowBuffer: framebuffer.NewShadow(1024, 1024),
		shadowShader: rendering.NewShader("shadow"),
	}
}

type Spot struct {
	BaseLight
	PointLight

	direction mgl32.Vec3
	// radians
	cutoff float32

	shadowBuffer *framebuffer.FBO
	shadowShader components.Shader
}

func (c *Spot) Direction() mgl32.Vec3 {
	r := c.Transform().TransformedRot()

	t := r.Rotate(mgl32.Vec3{0, 0, -1})
	return t
}

func (b *Spot) Cutoff() float32 {
	return b.cutoff
}

func (b *Spot) AddToEngine(e components.Engine) {
	e.GetRenderingEngine().AddLight(b)
}

func (b *Spot) ViewProjection() mgl32.Mat4 {
	const nearPlane float32 = 0.1
	const farPlane float32 = 30
	lightProjection := mgl32.Ortho(-8, 8, -8, 8, nearPlane, farPlane)
	//lightProjection := mgl32.Perspective(1.2, float32(1024/1024), nearPlane, farPlane)
	lightView := mgl32.LookAt(b.Position().X(), b.Position().Y(), b.Position().Z(), 0, 0, 0, 0, 1, 0)
	return lightProjection.Mul4(lightView)
}

func (b *Spot) BindShadowBuffer() {
	b.shadowBuffer.Bind()
	b.shadowBuffer.Texture().SetViewPort()
	gl.Enable(gl.DEPTH_TEST)
	gl.Clear(gl.DEPTH_BUFFER_BIT)
}

func (b *Spot) ShadowShader() components.Shader {
	return b.shadowShader
}

func (b *Spot) BindShadowTexture(samplerSlot uint32, samplerName string) {
	gl.ActiveTexture(gl.TEXTURE0 + uint32(samplerSlot))
	b.Shader().SetUniform(samplerName, int32(samplerSlot))
	gl.BindTexture(gl.TEXTURE_2D, b.shadowBuffer.Texture().ID())
}
