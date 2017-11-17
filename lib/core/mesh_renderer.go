package core

func NewMeshRenderer(mesh *Mesh) *MeshRenderer {
	return &MeshRenderer{
		mesh: mesh,
	}
}

type MeshRenderer struct {
	mesh *Mesh
	// material *Material
}

func (m *MeshRenderer) Render(shader *Shader) {
	shader.Bind()
	//shader.UpdateUniforms(GetTransform(), m_material, renderingEngine);
	m.mesh.Draw()
}

func (m *MeshRenderer) Input(*Input) {}
func (m *MeshRenderer) Update()      {}
