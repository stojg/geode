package core

import (
	"math"
	"math/rand"
	"time"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/stojg/graphics/lib/buffers"
	"github.com/stojg/graphics/lib/components"
	"github.com/stojg/graphics/lib/debug"
	"github.com/stojg/graphics/lib/resources"
)

const MaxParticles = 25000
const InstanceDataLength = 17 // floats (MAT4) + transparancy

func NewParticleSystem(perSecond float64) *ParticleSystem {

	model, vao := NewParticleModel()

	o := NewGameObject(components.R_PARTICLE)
	o.SetModel(model)

	vbo := buffers.CreateEmptyFloatVBO(vao, InstanceDataLength*MaxParticles, gl.STREAM_DRAW)
	buffers.AddInstancedAttribute(vao, vbo, 1, 4, InstanceDataLength, 0)
	buffers.AddInstancedAttribute(vao, vbo, 2, 4, InstanceDataLength, 4)
	buffers.AddInstancedAttribute(vao, vbo, 3, 4, InstanceDataLength, 8)
	buffers.AddInstancedAttribute(vao, vbo, 4, 4, InstanceDataLength, 12)
	buffers.AddInstancedAttribute(vao, vbo, 5, 1, InstanceDataLength, 16)

	return &ParticleSystem{
		GameObject: *o,
		perSecond:  perSecond,
		vbo:        vbo,
		vao:        vao,
	}
}

type ParticleSystem struct {
	GameObject
	perSecond   float64
	reminder    float64
	particles   []*Particle
	vao, vbo    uint32
	timeElapsed float64
}

func (s *ParticleSystem) Update(elapsed time.Duration) {
	s.timeElapsed += elapsed.Seconds()
}

var instanceData = make([]float32, MaxParticles*InstanceDataLength, MaxParticles*InstanceDataLength)

// @todo sort particles from back to front to fix blending
func (s *ParticleSystem) Draw(camera components.Viewable, shader components.Shader, state components.RenderState) {

	secs := float32(s.timeElapsed)
	for i := len(s.particles) - 1; i >= 0; i-- {
		alive := s.particles[i].Update(secs)
		// @todo better deletion to minimise allocations
		if !alive {
			s.particles = append(s.particles[:i], s.particles[i+1:]...)
		}
	}

	s.reminder += s.timeElapsed * s.perSecond
	toCreate, reminder := math.Modf(s.reminder)
	s.reminder = reminder

	posX := s.Transform().Pos()[0]
	posY := s.Transform().Pos()[1]
	posZ := s.Transform().Pos()[2]

	x, z := rand.Float32()*512-512/2, rand.Float32()*512-512/2
	for i := 0; i < int(toCreate); i++ {
		if len(s.particles) < MaxParticles {
			s.addParticle([3]float32{posX + rand.Float32()*x, posY, posZ + rand.Float32()*z}, [3]float32{rand.Float32()*0.5 - 0.25, rand.Float32() * 5, rand.Float32()*0.5 - 0.25}, rand.Float32()*0.05+0.025, rand.Float32()*45, 0.01, rand.Float32()*9+1)
		}
	}

	s.timeElapsed = 0

	count := 0
	for i := range s.particles {
		x := s.particles[i].Transform(camera)
		copy(instanceData[count:], x[0:16])
		instanceData[count+16] = s.particles[i].Transparency
		count += 17
	}
	buffers.UpdateFloatVBO(s.vao, s.vbo, len(instanceData), instanceData, gl.STREAM_DRAW)

	gl.BindVertexArray(s.vao)
	debug.AddVertexBind()
	gl.DrawArraysInstanced(gl.TRIANGLE_STRIP, 0, 4, int32(len(s.particles)))
	debug.Drawcall()

}

func (m *ParticleSystem) addParticle(pos, vel [3]float32, scale, rotAngle, gravity, life float32) {
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

func (p *ParticleModel) Material() components.Material {
	return p.material
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
	debug.AddVertexBind()
	gl.EnableVertexAttribArray(0)
	gl.EnableVertexAttribArray(1)
	gl.EnableVertexAttribArray(2)
	gl.EnableVertexAttribArray(3)
	gl.EnableVertexAttribArray(4)
	gl.EnableVertexAttribArray(5)
}

func (p *ParticleMesh) Draw() {
	gl.DrawArraysInstanced(gl.TRIANGLE_STRIP, 0, 4, 1)
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
	debug.AddVertexBind()
}

func setupVAO() uint32 {

	quadVertices := []float32{
		-0.5, -0.5, 0.0,
		0.5, -0.5, 0.0,
		-0.5, 0.5, 0.0,
		0.5, 0.5, 0.0,
	}

	// create arrays
	var vao uint32
	gl.GenVertexArrays(1, &vao)

	// load data into vertex buffer
	vbo := buffers.CreateFloatVBO(vao, len(quadVertices), quadVertices, gl.STATIC_DRAW)
	buffers.AddAttribute(vao, vbo, 0, 3, 3, 0)
	return vao
}
