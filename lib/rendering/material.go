package rendering

import "github.com/stojg/graphics/lib/components"

func NewMaterial() *Material {
	return &Material{
		textures: make(map[string]*Texture),
	}
}

type Material struct {
	textures map[string]*Texture
}

func (m *Material) AddTexture(name string, texture *Texture) {
	m.textures[name] = texture
}

func (m *Material) Texture(name string) components.Texture {
	texture, ok := m.textures[name]
	if !ok {
		panic("could not find texture, should return test texture instead")
	}
	return texture
}
