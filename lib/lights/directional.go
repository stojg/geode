package lights

import (
	"fmt"
	"math"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/stojg/graphics/lib/components"
)

func NewDirectional(shadowSize int, r, g, b, intensity float32) *Directional {
	var halfSize float32 = 10 / 2

	light := &Directional{
		BaseLight: BaseLight{
			color:       mgl32.Vec3{r, g, b}.Mul(intensity),
			constant:    0,
			maxDistance: 10,
			intensity:   intensity,
		},
	}

	if shadowSize != 0 {
		projection := mgl32.Ortho(-halfSize, halfSize, -halfSize, halfSize, -halfSize, halfSize)
		light.shadowInfo = NewShadowInfo(shadowSize, projection, false)
		light.shadowInfo.halfSize = halfSize
		light.shadowInfo.shadowVarianceMin = 0.00002
		light.shadowInfo.lightBleedReduction = 0.8
	}

	return light
}

type Directional struct {
	BaseLight
	matrix mgl32.Mat4
}

func (b *Directional) AddToEngine(e components.Engine) {
	e.RenderingEngine().AddLight(b)
}

func (b *Directional) SetCamera(inPos mgl32.Vec3, inRot mgl32.Quat) {

	fmt.Println(inPos)
	resultPos := inPos.Mul(-1).Add(inRot.Rotate(mgl32.Vec3{0, 0, b.shadowInfo.halfSize}))
	resultRot := b.Transform().TransformedRot().Conjugate()

	size := 1 << uint(b.shadowInfo.SizeAsPowerOfTwo())
	lightSpaceCameraPos := resultRot.Conjugate().Rotate(resultPos)
	texelSize := 2 * b.shadowInfo.halfSize / float32(size)
	lightSpaceCameraPos[0] = texelSize * float32(math.Floor(float64(lightSpaceCameraPos[0]))) / texelSize
	lightSpaceCameraPos[1] = texelSize * float32(math.Floor(float64(lightSpaceCameraPos[1]))) / texelSize
	lightSpaceCameraPos[2] = texelSize * float32(math.Floor(float64(lightSpaceCameraPos[2]))) / texelSize
	resultPos = resultRot.Rotate(lightSpaceCameraPos)

	translation := mgl32.Translate3D(resultPos[0], resultPos[1], resultPos[2])
	b.matrix = resultRot.Mat4().Mul4(translation)
}

func (b *Directional) Direction() mgl32.Vec3 {
	return b.Parent().Transform().TransformedPos().Normalize()
}

func (b *Directional) GetView() mgl32.Mat4 {
	//This comes from the conjugate rotation because the world should appear to rotate opposite to the camera's rotation.
	rotation := b.Transform().TransformedRot().Conjugate().Mat4()
	//Similarly, the translation is inverted because the world appears to move opposite to the camera's movement.
	position := b.Transform().TransformedPos().Mul(-1)
	translation := mgl32.Translate3D(position[0], position[1], position[2])
	return rotation.Mul4(translation)
}

func (b *Directional) ViewProjection() mgl32.Mat4 {
	return b.Projection().Mul4(b.matrix)
}
