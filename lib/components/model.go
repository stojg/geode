package components

func NewModel(mesh Drawable, material Material) *Model {
	return &Model{
		mesh:     mesh,
		material: material,
	}
}

type Model struct {
	GameComponent

	mesh     Drawable
	material Material
}

func (m *Model) Bind(shader Shader, engine RenderingEngine) {
	shader.UpdateUniforms(m.material, engine)
	m.mesh.Bind()
}

func (m *Model) Draw() {
	m.mesh.Draw()
}

func (m *Model) Unbind() {
	m.mesh.Unbind()
}
