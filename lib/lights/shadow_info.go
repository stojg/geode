package lights

import "github.com/go-gl/mathgl/mgl32"

func NewShadowInfo(sizeAsPowerOfTwo int, projection mgl32.Mat4, flipFaces bool) *ShadowInfo {
	return &ShadowInfo{
		sizeAsPowerOfTwo:    sizeAsPowerOfTwo,
		projection:          projection,
		flipFaces:           flipFaces,
		shadowVarianceMin:   0.00002,
		lightBleedReduction: 0.2,
	}
}

type ShadowInfo struct {
	sizeAsPowerOfTwo int
	projection       mgl32.Mat4
	flipFaces        bool

	shadowVarianceMin   float32
	lightBleedReduction float32
}

func (s *ShadowInfo) SizeAsPowerOfTwo() int {
	return s.sizeAsPowerOfTwo
}

func (s *ShadowInfo) LightBleedReduction() float32 {
	return s.lightBleedReduction
}

func (s *ShadowInfo) ShadowVarianceMin() float32 {
	return s.shadowVarianceMin
}

func (s *ShadowInfo) FlipFaces() bool {
	return s.flipFaces
}

func (s *ShadowInfo) Projection() mgl32.Mat4 {
	return s.projection
}
