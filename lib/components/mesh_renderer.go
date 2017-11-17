package components

func NewMeshRenderer(mesh Drawable) *MeshRenderer {
	return &MeshRenderer{
		mesh: mesh,
	}
}

type MeshRenderer struct {
	BaseComponent

	mesh Drawable
	// material *Material
}

func (m *MeshRenderer) Render(shader Bindable, engine UniformUpdater) {
	shader.Bind()
	//shader.UpdateUniforms(GetTransform(), m_material, renderingEngine);
	m.mesh.Draw()
}

