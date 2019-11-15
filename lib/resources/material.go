package resources

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/stojg/geode/lib/components"
)

func NewMaterial() *Material {
	return &Material{
		textures: make(map[string]components.Texture),
		vectors:  make(map[string]mgl32.Vec3),
		floats:   make(map[string]float32),
	}
}

type Material struct {
	textures map[string]components.Texture
	vectors  map[string]mgl32.Vec3
	floats   map[string]float32
}

func (m *Material) Float(n string) float32 {
	float, ok := m.floats[n]
	if ok {
		return float
	}
	return 0
}

func (m *Material) AddFloat(n string, a float32) {
	m.floats[n] = a
}

func (m *Material) Vector(n string) mgl32.Vec3 {
	vector, ok := m.vectors[n]
	if ok {
		return vector
	}
	return mgl32.Vec3{0, 0, 0}
}

func (m *Material) AddVector(n string, vec3 mgl32.Vec3) {
	m.vectors[n] = vec3
}

func (m *Material) Textures() map[string]components.Texture {
	return m.textures
}
func (m *Material) Texture(n string) components.Texture {
	texture, ok := m.textures[n]
	if ok {
		return texture
	}
	return m.textures[n]
}

func (m *Material) AddTexture(n string, t components.Texture) {
	m.textures[n] = t
}
