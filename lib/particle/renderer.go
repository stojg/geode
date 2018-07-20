package particle

import (
	"unsafe"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/stojg/graphics/lib/components"
	"github.com/stojg/graphics/lib/rendering/primitives"
	"github.com/stojg/graphics/lib/rendering/shader"
)

func NewRenderer(s components.RenderState) *Renderer {
	r := &Renderer{
		RenderState: s,
		shader:      shader.NewShader("particle"),
		quadVao:     setupVAO(),
	}

	//objVert, objInd, err := loader.Load("res/meshes/cube/model.obj")
	//if err != nil {
	//	fmt.Printf("Model loading failed: %v", err)
	//	os.Exit(1)
	//}
	//var meshes []*resources.Mesh
	//for i, data := range objVert {
	//	mesh := resources.NewMesh()
	//	mesh.SetVertices(resources.ConvertToVertices(data, objInd[i]), objInd[i])
	//	meshes = append(meshes, mesh)
	//}

	return r
}

type Renderer struct {
	components.RenderState
	shader  components.Shader
	quadVao uint32
	mesh    components.Drawable
}

func (r *Renderer) Render(objects components.Renderable) {
	r.shader.Bind()
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.DepthMask(false)
	//gl.BindVertexArray(r.quadVao)
	//gl.EnableVertexAttribArray(0)

	objects.Render(r.RenderState.Camera(), r.shader, r.RenderState, components.R_PARTICLE)

	//for _, p := range objects {
	//	r.shader.UpdateUniform("model", p.Transform(r.Camera()))
	//	r.shader.UpdateUniform("transparency", p.Transparency)
	//	gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, gl.PtrOffset(0))
	//	debug.Drawcall()
	//}

	//gl.DisableVertexAttribArray(0)
	//gl.BindVertexArray(0)
	gl.DepthMask(true)
	gl.Disable(gl.BLEND)
}

const sizeOfUint32 = unsafe.Sizeof(uint32(0))

func setupVAO() uint32 {
	var quadVao uint32
	quadVertices := []float32{
		-0.5, 0.5, 0, // top left
		-0.5, -0.5, 0, // bottom left
		0.5, 0.5, 0, // top right
		0.5, -0.5, 0, // bottom right
	}
	indices := []uint32{0, 1, 2, 1, 3, 2}

	// Create buffers/arrays
	var vbo, ebo uint32
	gl.GenBuffers(1, &vbo)
	gl.GenVertexArrays(1, &quadVao)
	gl.GenBuffers(1, &ebo)

	gl.BindVertexArray(quadVao)

	// load data into vertex buffer
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(quadVertices)*primitives.SizeOfFloat32*3, gl.Ptr(quadVertices), gl.STATIC_DRAW)

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*int(sizeOfUint32), gl.Ptr(indices), gl.STATIC_DRAW)

	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, int32(primitives.SizeOfFloat32)*3, gl.PtrOffset(0))

	// reset the current bound vertex array so that no one else mistakenly changes the VAO
	gl.BindVertexArray(0)

	return quadVao
}
