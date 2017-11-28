package lights

import (
	"fmt"
	"math"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/stojg/graphics/lib/components"
	"github.com/stojg/graphics/lib/rendering"
)

func NewPoint(r, g, b, intensity float32) *PointLight {
	color := mgl32.Vec3{r, g, b}
	pointLight := &PointLight{
		BaseLight: BaseLight{
			color:       color.Mul(intensity),
			shader:      rendering.NewShader("forward_point"),
			maxDistance: 3,
		},
		constant: 1,
		linear:   0.7,
		exponent: 1.80,
	}

	max := color[0]
	if color[1] > max {
		max = color[1]
	}
	if color[2] > max {
		max = color[2]
	}

	fmt.Println(color)
	fmt.Println(max)

	const colorDepth = 256

	{
		a := pointLight.Exponent()
		b := pointLight.Linear()
		c := pointLight.Constant() - colorDepth*intensity*max
		dist := (-b + float32(math.Sqrt(float64(b*b-4*a*c)))) / (2 * a)
		fmt.Println(dist)
		pointLight.BaseLight.maxDistance = dist
	}

	return pointLight
}

type PointLight struct {
	BaseLight
	constant float32
	linear   float32
	exponent float32
}

func (point *PointLight) AddToEngine(e components.Engine) {
	e.RenderingEngine().AddLight(point)
}

func (point *PointLight) Exponent() float32 {
	return point.exponent
}

func (point *PointLight) Linear() float32 {
	return point.linear
}

func (point *PointLight) Constant() float32 {
	return point.constant
}

/* @todo: calculate light range
float a = m_attenuation.Exponent();
float b = m_attenuation.Linear();
float c = m_attenuation.Constant() - COLOR_DEPTH * intensity * color.Max();
m_range = (-b + sqrtf(b*b - 4*a*c))/(2*a);
*/
