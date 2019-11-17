package resources

import (
	"math"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/stojg/geode/lib/buffers"
	"github.com/stojg/geode/lib/components"
	"github.com/stojg/geode/lib/debug"
	"github.com/stojg/geode/lib/geometry"
)

// @todo check http://ogldev.atspace.co.uk/www/tutorial33/tutorial33.html for proper instanced rendering
// @todo Would be nice to have Model that has one or many meshes and textures
func NewMesh() *Mesh {
	m := &Mesh{}

	return m
}

type Mesh struct {
	vbo  uint32
	vao  uint32
	ebo  uint32
	num  int32
	aabb *geometry.AABB
}

func (m *Mesh) SetVertices(vertices []Vertex, indices []uint32) {
	// Create buffers/arrays
	gl.GenVertexArrays(1, &m.vao)
	gl.BindVertexArray(m.vao)
	debug.AddVertexBind()

	m.num = int32(len(indices))

	buffers.CreateIntEBO(m.vao, len(indices), indices, gl.STATIC_DRAW)

	// load data into vertex buffers
	m.vbo = buffers.CreateVBO(m.vao, len(vertices)*int(sizeOfVertex), vertices, gl.STATIC_DRAW)

	// position
	buffers.AddAttribute(m.vao, m.vbo, 0, 3, 11, 0)
	// normals
	buffers.AddAttribute(m.vao, m.vbo, 1, 3, 11, 3)
	// texture coordinates
	buffers.AddAttribute(m.vao, m.vbo, 2, 2, 11, 6)
	// tangents
	buffers.AddAttribute(m.vao, m.vbo, 3, 3, 11, 8)

	m.aabb = &geometry.AABB{}

	center := mgl32.Vec3{}
	extent := mgl32.Vec3{}

	center[0], extent[0] = halfWidth(vertices, [3]float32{1, 0, 0})
	center[1], extent[1] = halfWidth(vertices, [3]float32{0, 1, 0})
	center[2], extent[2] = halfWidth(vertices, [3]float32{0, 0, 1})
	m.aabb.SetC(center)
	m.aabb.SetR(extent)
}

func (m *Mesh) AABB() components.AABB {
	return m.aabb
}

func (m *Mesh) Bind() {
	gl.BindVertexArray(m.vao)
	debug.AddVertexBind()
	gl.EnableVertexAttribArray(0)
	gl.EnableVertexAttribArray(1)
	gl.EnableVertexAttribArray(2)
	gl.EnableVertexAttribArray(3)
}

func (m *Mesh) Draw() {
	gl.DrawElements(gl.TRIANGLES, m.num, gl.UNSIGNED_INT, gl.PtrOffset(0))
	debug.Drawcall()
}

func (m *Mesh) Unbind() {
	gl.DisableVertexAttribArray(3)
	gl.DisableVertexAttribArray(2)
	gl.DisableVertexAttribArray(1)
	gl.DisableVertexAttribArray(0)
	gl.BindVertexArray(0)
	debug.AddVertexBind()
}

func (m *Mesh) CleanUp() {
	gl.DeleteBuffers(1, &m.vbo)
	gl.DeleteBuffers(1, &m.ebo)
	gl.DeleteVertexArrays(1, &m.vao)
}

func halfWidth(in []Vertex, direction [3]float32) (float32, float32) {
	min, max := float32(math.MaxFloat32), float32(-math.MaxFloat32)
	var proj float32
	for i := 0; i < len(in); i++ {
		dot(in[i].Pos, direction, &proj)
		if proj < min {
			min = proj
		}
		if proj > max {
			max = proj
		}
	}

	center := 0.5 * (max + min)
	extent := 0.5 * (max - min)
	return center, extent
}

func CalcMinMax(in []Vertex, direction [3]float32) [2]float32 {
	min, max := float32(math.MaxFloat32), float32(-math.MaxFloat32)
	var proj float32
	for i := 0; i < len(in); i++ {
		dot(in[i].Pos, direction, &proj)
		if proj < min {
			min = proj
		}
		if proj > max {
			max = proj
		}
	}

	return [2]float32{min, max}
}

func dot(a [3]float32, b [3]float32, result *float32) {
	x := a[0] * b[0]
	y := a[1] * b[1]
	z := a[2] * b[2]
	*result = x + y + z
}
