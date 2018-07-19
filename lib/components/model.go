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

func (m *Model) Bind(shader Shader, engine RenderState) {
	shader.UpdateUniforms(m.material, engine)
	m.mesh.Bind()
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
	return m.mesh.HalfWidths()
}
