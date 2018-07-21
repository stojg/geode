package core

import (
	"math"
	"math/rand"
	"time"
	"unsafe"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/stojg/graphics/lib/components"
	"github.com/stojg/graphics/lib/debug"
	"github.com/stojg/graphics/lib/rendering/primitives"
	"github.com/stojg/graphics/lib/resources"
)

const MaxParticles = 10000
const InstanceDataLength = 17 // floats (MAT4) + transparancy

func NewParticleSystem(perSecond float64) *ParticleSystem {

	model, vao := NewParticleModel()

	o := NewGameObject(components.R_PARTICLE)
	o.SetModel(model)

	vbo := createEmptyFloatBO(InstanceDataLength*MaxParticles, gl.ARRAY_BUFFER, gl.STREAM_DRAW)
	addInstancedAttribute(vao, vbo, 1, 4, InstanceDataLength, 0)
	addInstancedAttribute(vao, vbo, 2, 4, InstanceDataLength, 4)
	addInstancedAttribute(vao, vbo, 3, 4, InstanceDataLength, 8)
	addInstancedAttribute(vao, vbo, 4, 4, InstanceDataLength, 12)
	addInstancedAttribute(vao, vbo, 5, 1, InstanceDataLength, 16)

	return &ParticleSystem{
		GameObject:     *o,
		perSecond:      perSecond,
		perInstanceVBO: vbo,
	}
}

type ParticleSystem struct {
	GameObject
	perSecond      float64
	reminder       float64
	particles      []*Particle
	perInstanceVBO uint32
}

func (s *ParticleSystem) Update(elapsed time.Duration) {
	secs := float32(elapsed.Seconds())
	for i := len(s.particles) - 1; i >= 0; i-- {
		alive := s.particles[i].Update(secs)
		// @todo better deletion to minimise allocations
		if !alive {
			s.particles = append(s.particles[:i], s.particles[i+1:]...)
		}
	}

	s.reminder += elapsed.Seconds() * s.perSecond
	toCreate, reminder := math.Modf(s.reminder)
	s.reminder = reminder
	for i := 0; i < int(toCreate); i++ {
		s.AddParticle(s.Transform().Pos(), [3]float32{rand.Float32()*4 - 2, rand.Float32()*15 + 5, rand.Float32()*4 - 2}, rand.Float32()*0.05+0.025, rand.Float32()*45, 1, rand.Float32()*9+1)
	}
}

func (s *ParticleSystem) OldDraw(camera components.Viewable, shader components.Shader, state components.RenderState) {
	for _, p := range s.particles {
		shader.UpdateUniform("model", p.Transform(camera))
		s.model.Draw()
	}
}

// @todo sort particles from back to front to fix blending
func (s *ParticleSystem) Draw(camera components.Viewable, shader components.Shader, state components.RenderState) {

	var instanceData []float32

	for _, p := range s.particles {
		x := p.Transform(camera)
		for i := 0; i < 4; i++ {
			for _, j := range x.Col(i) {
				instanceData = append(instanceData, j)
			}
		}
		instanceData = append(instanceData, p.Transparency)
	}
	updateVBO(s.perInstanceVBO, instanceData, gl.ARRAY_BUFFER, gl.STREAM_DRAW)
	gl.DrawElementsInstanced(gl.TRIANGLES, 6, gl.UNSIGNED_INT, gl.PtrOffset(0), int32(len(s.particles)))
	debug.Drawcall()

}

func (m *ParticleSystem) AddParticle(pos, vel [3]float32, scale, rotAngle, gravity, life float32) {
	p := NewParticle(pos, vel, scale, rotAngle, gravity, life)
	m.particles = append(m.particles, p)
}

func (s *ParticleSystem) IsVisible(camera components.Viewable) bool {
	return true
}

func NewParticleModel() (components.Model, uint32) {

	vao := setupVAO()
	return &ParticleModel{
		mesh:     &ParticleMesh{vao: vao},
		material: resources.NewMaterial(),
	}, vao
}

type ParticleModel struct {
	mesh     components.Drawable
	material components.Material
}

func (p *ParticleModel) AABB() [3][2]float32 {
	return p.mesh.AABB()
}

func (p *ParticleModel) Bind(shader components.Shader, state components.RenderState) {
	shader.UpdateUniforms(p.material, state)
	p.mesh.Bind()
}

func (p *ParticleModel) Draw() {
	p.mesh.Draw()
}

func (p *ParticleModel) Unbind() {
	p.mesh.Unbind()
}

type ParticleMesh struct {
	vao uint32
	vpb uint32
}

func (p *ParticleMesh) AABB() [3][2]float32 {
	return [3][2]float32{}
}

func (p *ParticleMesh) Bind() {
	gl.BindVertexArray(p.vao)
	gl.EnableVertexAttribArray(0)
	gl.EnableVertexAttribArray(1)
	gl.EnableVertexAttribArray(2)
	gl.EnableVertexAttribArray(3)
	gl.EnableVertexAttribArray(4)
	gl.EnableVertexAttribArray(5)
}

func (p *ParticleMesh) Draw() {
	gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, gl.PtrOffset(0))
	debug.Drawcall()
}

func (p *ParticleMesh) Unbind() {
	gl.DisableVertexAttribArray(5)
	gl.DisableVertexAttribArray(4)
	gl.DisableVertexAttribArray(3)
	gl.DisableVertexAttribArray(2)
	gl.DisableVertexAttribArray(1)
	gl.DisableVertexAttribArray(0)
	gl.BindVertexArray(0)
}

const sizeOfUint32 = unsafe.Sizeof(uint32(0))

func setupVAO() uint32 {

	quadVertices := []float32{
		-0.5, 0.5, 0, // top left
		-0.5, -0.5, 0, // bottom left
		0.5, 0.5, 0, // top right
		0.5, -0.5, 0, // bottom right
	}
	indices := []uint32{0, 1, 2, 1, 3, 2}

	// create arrays
	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	// load data into vertex buffer
	vbo := createEmptyFloatBO(len(quadVertices), gl.ARRAY_BUFFER, gl.STATIC_DRAW)
	updateVBO(vbo, quadVertices, gl.ARRAY_BUFFER, gl.STATIC_DRAW)
	addAttribute(vao, vbo, gl.ARRAY_BUFFER, 0, 3, 3, 0)

	// Create buffers
	_ = createUint32BO(vao, gl.ELEMENT_ARRAY_BUFFER, gl.STATIC_DRAW, indices)

	return vao
}

func createUint32BO(vao uint32, target, usage uint32, indices []uint32) uint32 {
	var ebo uint32
	gl.GenBuffers(1, &ebo)
	gl.BindVertexArray(vao)
	gl.BindBuffer(target, ebo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*int(sizeOfUint32), gl.Ptr(indices), usage)
	gl.BindVertexArray(0)
	return ebo
}

func createEmptyFloatBO(floatCount int, target, usage uint32) uint32 {
	var bufferObject uint32
	gl.GenBuffers(1, &bufferObject)
	gl.BindBuffer(target, bufferObject)
	gl.BufferData(target, floatCount*primitives.SizeOfFloat32, nil, usage)
	gl.BindBuffer(target, 0)
	return bufferObject
}

func addInstancedAttribute(vao, vbo uint32, attribute uint32, dataSizeInFloats int32, instanceDataLength int, offset int) {
	gl.BindVertexArray(vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.VertexAttribPointer(attribute, dataSizeInFloats, gl.FLOAT, false, int32(instanceDataLength*primitives.SizeOfFloat32), gl.PtrOffset(offset*primitives.SizeOfFloat32))
	gl.VertexAttribDivisor(attribute, 1)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindVertexArray(0)
}

func addAttribute(vao, vbo uint32, target, attribute uint32, dataSizeInFloats int32, instanceDataLength int, offset int) {
	gl.BindVertexArray(vao)
	gl.BindBuffer(target, vbo)
	gl.VertexAttribPointer(attribute, dataSizeInFloats, gl.FLOAT, false, int32(instanceDataLength*primitives.SizeOfFloat32), gl.PtrOffset(offset*primitives.SizeOfFloat32))
	gl.BindBuffer(target, 0)
	gl.BindVertexArray(0)
}

func updateVBO(vbo uint32, data []float32, target, usage uint32) {
	gl.BindBuffer(target, vbo)
	// Buffer orphaning, a common way to improve streaming perf.
	gl.BufferData(target, len(data)*primitives.SizeOfFloat32, nil, usage)
	gl.BufferSubData(target, 0, len(data)*primitives.SizeOfFloat32, gl.Ptr(data))
	gl.BindBuffer(target, 0)

}
