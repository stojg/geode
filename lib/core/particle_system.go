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

const MaxParticles = 500000
const InstanceDataLength = 4

type particleData struct {
	aliveCount   int
	velocity     [MaxParticles][3]float32
	gravity      [MaxParticles]float32
	position     [MaxParticles][3]float32
	transparency [MaxParticles]float32
	rotation     [MaxParticles]float32
	scale        [MaxParticles]float32
	elapsedTime  [MaxParticles]float32
	alive        [MaxParticles]bool
	lifeLength   [MaxParticles]float32
}

func (p *particleData) add(pos, vel [3]float32, scale, rotAngle, gravity, life float32) {
	if p.aliveCount >= MaxParticles {
		return
	}
	p.alive[p.aliveCount] = true
	p.transparency[p.aliveCount] = 1.0
	p.position[p.aliveCount] = pos
	p.velocity[p.aliveCount] = vel
	p.scale[p.aliveCount] = scale
	p.rotation[p.aliveCount] = rotAngle
	p.gravity[p.aliveCount] = gravity * Gravity
	p.lifeLength[p.aliveCount] = life
	p.elapsedTime[p.aliveCount] = 0

	p.aliveCount++
}

func (p *particleData) remove(id int) {
	p.alive[id] = false
	p.swap(id, p.aliveCount-1)
	p.aliveCount--
}

func (p *particleData) swap(a, b int) {
	p.velocity[a], p.velocity[b] = p.velocity[b], p.velocity[a]
	p.gravity[a], p.gravity[b] = p.gravity[b], p.gravity[a]
	p.position[a], p.position[b] = p.position[b], p.position[a]
	p.transparency[a], p.transparency[b] = p.transparency[b], p.transparency[a]
	p.rotation[a], p.rotation[b] = p.rotation[b], p.rotation[a]
	p.scale[a], p.scale[b] = p.scale[b], p.scale[a]
	p.elapsedTime[a], p.elapsedTime[b] = p.elapsedTime[b], p.elapsedTime[a]
	p.alive[a], p.alive[b] = p.alive[b], p.alive[a]
	p.lifeLength[a], p.lifeLength[b] = p.lifeLength[b], p.lifeLength[a]
}

func NewParticleSystem(perSecond float64) *ParticleSystem {

	model, vao := NewParticleModel()

	o := NewGameObject(components.R_PARTICLE)
	o.SetModel(model)

	vbo := buffers.CreateEmptyFloatVBO(vao, InstanceDataLength*MaxParticles, gl.STREAM_DRAW)
	buffers.AddInstancedAttribute(vao, vbo, 1, 4, InstanceDataLength, 0)

	return &ParticleSystem{
		GameObject: *o,
		perSecond:  perSecond,
		vbo:        vbo,
		vao:        vao,
		data:       &particleData{},
	}
}

func simpleUpdater(data *particleData, elapsed float32) {
	for i := 0; i < data.aliveCount; i++ {
		data.velocity[i][1] += data.gravity[i] * elapsed
	}

	for i := 0; i < data.aliveCount; i++ {
		data.position[i][0] += data.velocity[i][0] * elapsed
		data.position[i][1] += data.velocity[i][1] * elapsed
		data.position[i][2] += data.velocity[i][2] * elapsed
	}

	for i := 0; i < data.aliveCount; i++ {
		data.elapsedTime[i] += elapsed
	}

	for i := 0; i < data.aliveCount; i++ {
		data.transparency[i] = 1 - data.elapsedTime[i]/data.lifeLength[i]
	}

	for i := 0; i < data.aliveCount; i++ {
		alive := data.elapsedTime[i] < data.lifeLength[i]
		if !alive {
			data.remove(i)
		}
	}
}

type ParticleSystem struct {
	GameObject
	perSecond   float64
	reminder    float64
	vao, vbo    uint32
	timeElapsed float64

	data *particleData
}

func (s *ParticleSystem) Update(elapsed time.Duration) {
	s.timeElapsed += elapsed.Seconds()
}

var instanceData = make([]float32, MaxParticles*InstanceDataLength, MaxParticles*InstanceDataLength)

// @todo sort particles from back to front to fix blending
func (s *ParticleSystem) Draw(camera components.Viewable, shader components.Shader, state components.RenderState) {

	elapsed := float32(s.timeElapsed)
	toCreate := s.calculateToCreate()
	//x, z := rand.Float32()*0, rand.Float32()*0
	posX := s.Transform().Pos()[0]
	posY := s.Transform().Pos()[1]
	posZ := s.Transform().Pos()[2]
	for i := 0; i < int(toCreate); i++ {
		s.data.add([3]float32{posX, posY, posZ}, [3]float32{rand.Float32()*0.5 - 0.25, rand.Float32() * 15, rand.Float32()*0.5 - 0.25}, rand.Float32()*0.05+0.025, 0, 1, rand.Float32()*4+1)
	}

	simpleUpdater(s.data, elapsed)

	s.updateInstanceData()
	buffers.UpdateFloatVBO(s.vao, s.vbo, len(instanceData), instanceData, gl.STREAM_DRAW)

	view := camera.View()
	s.state.SetVector3f("x_camRight", [3]float32{view[0], view[4], view[8]})
	s.state.SetVector3f("x_camUp", [3]float32{view[1], view[5], view[9]})

	gl.BindVertexArray(s.vao)
	debug.AddVertexBind()
	gl.DrawArraysInstanced(gl.TRIANGLE_STRIP, 0, 4, int32(s.data.aliveCount))
	debug.Drawcall()
	debug.SetParticles(uint64(s.data.aliveCount))

}

func (s *ParticleSystem) calculateToCreate() float64 {
	var toCreate float64
	s.reminder += s.timeElapsed * s.perSecond
	s.timeElapsed = 0
	toCreate, s.reminder = math.Modf(s.reminder)
	return toCreate
}

func (s *ParticleSystem) updateInstanceData() {
	count := 0
	for i := 0; i < s.data.aliveCount; i++ {
		copy(instanceData[count:], s.data.position[i][0:3])
		instanceData[count+InstanceDataLength-1] = s.data.transparency[i]
		count += InstanceDataLength
	}
}

func (s *ParticleSystem) IsVisible(camera components.Viewable) bool {
	return true
}

func NewParticleModel() (components.Model, uint32) {

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
}

func (p *ParticleMesh) Draw() {
	gl.DrawArraysInstanced(gl.TRIANGLE_STRIP, 0, 4, 1)
	debug.Drawcall()
}

func (p *ParticleMesh) Unbind() {
	gl.DisableVertexAttribArray(2)
	gl.DisableVertexAttribArray(1)
	gl.BindVertexArray(0)
	debug.AddVertexBind()
}
