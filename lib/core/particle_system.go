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

func NewParticleSystem(perSecond float64) *ParticleSystem {

	model := NewParticleModel()

	o := NewGameObject(components.R_PARTICLE)
	o.SetModel(model)

	return &ParticleSystem{
		GameObject: *o,
		perSecond:  perSecond,
	}
}

type ParticleSystem struct {
	GameObject
	perSecond float64
	reminder  float64
	particles []*Particle
}

func (s *ParticleSystem) Update(elapsed time.Duration) {
	secs := float32(elapsed.Seconds())
	for i := len(s.particles) - 1; i >= 0; i-- {
		alive := s.particles[i].Update(secs)
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

func (s *ParticleSystem) Draw(camera components.Viewable, shader components.Shader, state components.RenderState) {
	for _, p := range s.particles {
		shader.UpdateUniform("model", p.Transform(camera))
		shader.UpdateUniform("transparency", p.Transparency)
		s.model.Draw()
	}
}

func (m *ParticleSystem) AddParticle(pos, vel [3]float32, scale, rotAngle, gravity, life float32) {
	p := NewParticle(pos, vel, scale, rotAngle, gravity, life)
	m.particles = append(m.particles, p)
}

func (s *ParticleSystem) IsVisible(camera components.Viewable) bool {
	return true
}

func NewParticleModel() components.Model {
	mtrl := resources.NewMaterial()
	return &ParticleModel{
		mesh:     &ParticleMesh{vao: setupVAO()},
		material: mtrl,
	}
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
	vbo uint32
	vao uint32
	ebo uint32
	num int32
}

func (p *ParticleMesh) AABB() [3][2]float32 {
	return [3][2]float32{}
}

func (p *ParticleMesh) Bind() {
	gl.BindVertexArray(p.vao)
	gl.EnableVertexAttribArray(0)
}

func (p *ParticleMesh) Draw() {
	gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, gl.PtrOffset(0))
	debug.Drawcall()
}

func (p *ParticleMesh) Unbind() {
	gl.DisableVertexAttribArray(0)
	gl.BindVertexArray(0)
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
