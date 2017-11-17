package components

func NewMeshRenderer(mesh Drawable) *MeshRenderer {
	return &MeshRenderer{
		mesh: mesh,
	}
}

type MeshRenderer struct {
	GameComponent

	mesh Drawable
	// material *Material
}

func (m *MeshRenderer) Render(shader Shader, engine RenderingEngine) {
	shader.Bind()
	shader.UpdateUniforms(m.Transform(), engine)
	m.mesh.Draw()
}
