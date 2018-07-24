package core

import "github.com/stojg/graphics/lib/components"

func NewModel(mesh components.Drawable, material components.Material) *Model {
	return &Model{
		mesh:         mesh,
		material:     material,
		numberOfRows: 1,
	}
}

type Model struct {
	components.GameComponent
	mesh         components.Drawable
	material     components.Material
	numberOfRows uint32
}

func (m *Model) Bind(shader components.Shader, engine components.RenderState) {
	shader.UpdateUniforms(m.material, engine)
	m.mesh.Bind()
}

func (m *Model) Material() components.Material {
	return m.material
}

func (m *Model) Draw() {
	m.mesh.Draw()
}

func (m *Model) Unbind() {
	m.mesh.Unbind()
}

func (m *Model) AABB() [3][2]float32 {
	if m.mesh == nil {
		panic("no mesh?!")
	}
	return m.mesh.AABB()
}
