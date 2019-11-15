package core

import (
	"time"

	"github.com/stojg/geode/lib/components"
)

func NewModel(mesh components.Drawable, material components.Material) *Model {
	return &Model{
		mesh:     mesh,
		material: material,
	}
}

type Model struct {
	mesh     components.Drawable
	material components.Material
}

func (m *Model) Bind(shader components.Shader, engine components.RenderState) {
	shader.UpdateUniforms(m.material, engine)
	m.mesh.Bind()
}

func (m *Model) Update(time.Duration) {}

func (m *Model) Material() components.Material {
	return m.material
}

func (m *Model) Draw() {
	m.mesh.Draw()
}

func (m *Model) Unbind() {
	m.mesh.Unbind()
}

func (m *Model) AABB() components.AABB {
	return m.mesh.AABB()
}
