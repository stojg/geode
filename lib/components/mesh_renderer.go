package components

type Texture interface {
	Bind(samplerSlot uint32)
	BindAsRenderTarget()
	Width() int32
	Height() int32
}

type Material interface {
	Texture(name string) Texture
}

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
