package resources

import (
	"fmt"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/stojg/graphics/lib/components"
)

func NewMaterial() *Material {
	return &Material{
		textures:  make(map[string]*Texture),
		albedo:    mgl32.Vec3{0.02, 0.02, 0.02}, // charcoal
		metallic:  0.00,                         // non metallic
		roughness: 0.5,
	}
}

// https://docs.unrealengine.com/latest/INT/Engine/Rendering/Materials/PhysicallyBased/
type Material struct {
	textures  map[string]*Texture
	albedo    mgl32.Vec3
	metallic  float32
	roughness float32
}

func (m *Material) Albedo() mgl32.Vec3 {
	return m.albedo
}

func (m *Material) SetAlbedo(albedo mgl32.Vec3) {
	m.albedo = albedo
}

func (m *Material) SetMetallic(metallic float32) {
	m.metallic = metallic
}

func (m *Material) Metallic() float32 {
	return m.metallic
}

func (m *Material) SetRoughness(roughness float32) {
	m.roughness = roughness
}

func (m *Material) Roughness() float32 {
	return m.roughness
}

func (m *Material) AddTexture(name string, texture *Texture) {
	m.textures[name] = texture
}

func (m *Material) Texture(name string) components.Texture {
	texture, ok := m.textures[name]
	if !ok {
		panic(fmt.Sprintf("could not find texture '%s', should return test texture instead\n", name))
	}
	return texture
}
