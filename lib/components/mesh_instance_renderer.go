package components

func NewMeshInstanceRenderer(mesh Drawable, material Material) *MeshInstanceRenderer {
	return &MeshInstanceRenderer{
		mesh:     mesh,
		material: material,
	}
}

type MeshInstanceRenderer struct {
	GameComponent

	mesh     Drawable
	material Material
}

func (m *MeshInstanceRenderer) Render(shader Shader, engine RenderingEngine) {
	shader.Bind()
	m.mesh.Prepare()
	shader.UpdateUniforms(m.material, engine)
	for _, t := range m.parent.AllTransforms() {
		shader.UpdateTransform(t, engine)
		m.mesh.Draw()
	}
}
