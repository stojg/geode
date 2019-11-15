package lights

import (
	"math"
	"time"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/stojg/geode/lib/components"
)

func NewDirectional(shadowSize int, r, g, b, intensity float32) *Directional {
	var halfSize float32 = 40 / 2

	light := &Directional{
		BaseLight: BaseLight{
			color:    mgl32.Vec3{r, g, b}.Mul(intensity),
			constant: 0,
		},
	}

	if shadowSize != 0 {
		projection := mgl32.Ortho(-halfSize, halfSize, -halfSize, halfSize, -halfSize, halfSize)
		light.shadowInfo = NewShadowInfo(shadowSize, projection, false)
		light.shadowInfo.halfSize = halfSize
		light.shadowInfo.shadowVarianceMin = 0.00002
		light.shadowInfo.lightBleedReduction = 0.5
	}

	return light
}

type Directional struct {
	BaseLight
	matrix        mgl32.Mat4
	viewDirection mgl32.Mat4
	direction     mgl32.Vec3
}

func (direction *Directional) Update(elapsed time.Duration) {
	direction.BaseLight.Update(elapsed)
	direction.viewDirection = direction.Projection().Mul4(direction.matrix)
	direction.direction = direction.Parent().Transform().TransformedPos().Normalize()
}

func (direction *Directional) SetCamera(inPos mgl32.Vec3, inRot mgl32.Quat) {
	resultPos := inPos.Mul(-1).Add(inRot.Rotate(mgl32.Vec3{0, 0, direction.shadowInfo.halfSize}))
	resultRot := direction.Transform().TransformedRot().Conjugate()

	size := 1 << uint(direction.shadowInfo.SizeAsPowerOfTwo())
	lightSpaceCameraPos := resultRot.Conjugate().Rotate(resultPos)
	texelSize := 2 * direction.shadowInfo.halfSize / float32(size)
	lightSpaceCameraPos[0] = texelSize * float32(math.Floor(float64(lightSpaceCameraPos[0]))) / texelSize
	lightSpaceCameraPos[1] = texelSize * float32(math.Floor(float64(lightSpaceCameraPos[1]))) / texelSize
	lightSpaceCameraPos[2] = texelSize * float32(math.Floor(float64(lightSpaceCameraPos[2]))) / texelSize
	resultPos = resultRot.Rotate(lightSpaceCameraPos)

	translation := mgl32.Translate3D(resultPos[0], resultPos[1], resultPos[2])
	direction.matrix = resultRot.Mat4().Mul4(translation)
}

func (direction *Directional) Direction() mgl32.Vec3 {
	return direction.direction
}

func (direction *Directional) AddToEngine(e components.RenderState) {
	e.AddLight(direction)
}

func (direction *Directional) ViewProjection() mgl32.Mat4 {
	return direction.viewDirection
}
