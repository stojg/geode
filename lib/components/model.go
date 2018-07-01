package components

func NewModel(mesh Drawable, material Material) *Model {
	return &Model{
		mesh:         mesh,
		material:     material,
		numberOfRows: 1,
	}
}

type Model struct {
	GameComponent

	mesh         Drawable
	material     Material
	numberOfRows uint32
}

func (m *Model) NumberOfRows() uint32 {
	return m.numberOfRows
}

func (m *Model) SetNumberOfRows(numberOfRows uint32) {
	m.numberOfRows = numberOfRows
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
