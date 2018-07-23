package core

import (
	"math"
	"math/rand"
	"time"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/stojg/graphics/lib/components"
	"github.com/stojg/graphics/lib/debug"
	"github.com/stojg/graphics/lib/resources"
	"github.com/stojg/graphics/lib/utilities"
)

const MaxParticles = 10000
const InstanceDataLength = 17 // floats (MAT4) + transparancy

func NewParticleSystem(perSecond float64) *ParticleSystem {

	model, vao := NewParticleModel()

	o := NewGameObject(components.R_PARTICLE)
	o.SetModel(model)

	vbo := utilities.CreateEmptyVBO(vao, InstanceDataLength*MaxParticles, gl.STREAM_DRAW)
	utilities.AddInstancedAttribute(vao, vbo, 1, 4, InstanceDataLength, 0)
	utilities.AddInstancedAttribute(vao, vbo, 2, 4, InstanceDataLength, 4)
	utilities.AddInstancedAttribute(vao, vbo, 3, 4, InstanceDataLength, 8)
	utilities.AddInstancedAttribute(vao, vbo, 4, 4, InstanceDataLength, 12)
	utilities.AddInstancedAttribute(vao, vbo, 5, 1, InstanceDataLength, 16)

	return &ParticleSystem{
		GameObject: *o,
		perSecond:  perSecond,
		vbo:        vbo,
		vao:        vao,
	}
}

type ParticleSystem struct {
	GameObject
	perSecond float64
	reminder  float64
	particles []*Particle
	vao, vbo  uint32
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

	posX := s.Transform().Pos()[0]
	posY := s.Transform().Pos()[1]
	posZ := s.Transform().Pos()[2]

	for i := 0; i < int(toCreate); i++ {
		s.addParticle([3]float32{posX + rand.Float32()*100 - 50, posY, posZ + rand.Float32()*100 - 50}, [3]float32{rand.Float32()*0.5 - 0.25, rand.Float32() * 5, rand.Float32()*0.5 - 0.25}, rand.Float32()*0.05+0.025, rand.Float32()*45, 0.01, rand.Float32()*9+1)
	}
}

var instanceData = make([]float32, MaxParticles*InstanceDataLength, MaxParticles*InstanceDataLength)

// @todo sort particles from back to front to fix blending
func (s *ParticleSystem) Draw(camera components.Viewable, shader components.Shader, state components.RenderState) {

	count := 0

	for _, p := range s.particles {
		x := p.Transform(camera)
		for i := 0; i < 4; i++ {
			for _, j := range x.Col(i) {
				instanceData[count] = j
				count++
			}
		}
		instanceData[count] = p.Transparency
		count++
	}
	utilities.UpdateFloatVBO(s.vao, s.vbo, len(instanceData), instanceData, gl.STREAM_DRAW)

	gl.BindVertexArray(s.vao)
	gl.DrawElementsInstanced(gl.TRIANGLES, 6, gl.UNSIGNED_INT, gl.PtrOffset(0), int32(len(s.particles)))
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

	// load data into vertex buffer
	vbo := utilities.CreateFloatVBO(vao, len(quadVertices), quadVertices, gl.STATIC_DRAW)
	utilities.AddAttribute(vao, vbo, 0, 3, 3, 0)

	utilities.CreateIntEBO(vao, len(indices), indices, gl.STATIC_DRAW)

	return vao
}
