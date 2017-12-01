package components

func NewMeshRenderer(mesh Drawable, material Material) *MeshRenderer {
	return &MeshRenderer{
		mesh:     mesh,
		material: material,
	}
}

type MeshRenderer struct {
	GameComponent

	mesh     Drawable
	material Material
}

func (m *MeshRenderer) Render(shader Shader, engine RenderingEngine) {
	shader.Bind()
	shader.UpdateUniforms(m.Transform(), m.material, engine)
	m.mesh.Draw()
}
