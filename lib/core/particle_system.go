package core

import (
	"math"
	"math/rand"
	"sort"
	"time"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/stojg/geode/lib/buffers"
	"github.com/stojg/geode/lib/components"
	"github.com/stojg/geode/lib/debug"
	"github.com/stojg/geode/lib/physics"
	"github.com/stojg/geode/lib/resources"
)

const MaxParticles = 500000
const InstanceDataLength = 8
const Gravity = -9.82

type particleData struct {
	aliveCount       int
	velocity         [MaxParticles][3]float32
	gravity          [MaxParticles]float32
	position         [MaxParticles][3]float32
	transparency     [MaxParticles]float32
	rotation         [MaxParticles]float32
	scale            [MaxParticles]float32
	elapsedTime      [MaxParticles]float32
	alive            [MaxParticles]bool
	lifeLength       [MaxParticles]float32
	distanceToCamera [MaxParticles]float32
	colour           [MaxParticles][3]float32
}

func (p *particleData) add(pos, vel, colour [3]float32, scale, rotAngle, gravity, life float32) {
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
	p.distanceToCamera[p.aliveCount] = 0
	p.colour[p.aliveCount] = colour
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
	p.distanceToCamera[a], p.distanceToCamera[b] = p.distanceToCamera[b], p.distanceToCamera[a]
	p.colour[a], p.colour[b] = p.colour[b], p.colour[a]
}

func (a *particleData) Len() int           { return a.aliveCount }
func (a *particleData) Swap(i, j int)      { a.swap(i, j) }
func (a *particleData) Less(i, j int) bool { return a.distanceToCamera[i] > a.distanceToCamera[j] }

func NewParticleSystem(perSecond float64) *ParticleSystem {

	model, vao := NewParticleModel()

	o := NewGameObject(components.R_PARTICLE)
	o.SetModel(model)

	vbo := buffers.CreateEmptyFloatVBO(vao, InstanceDataLength*MaxParticles, gl.STREAM_DRAW)
	buffers.AddInstancedAttribute(vao, vbo, 1, 4, InstanceDataLength, 0)
	buffers.AddInstancedAttribute(vao, vbo, 2, 4, InstanceDataLength, 4)

	return &ParticleSystem{
		GameObject: *o,
		perSecond:  perSecond,
		vbo:        vbo,
		vao:        vao,
		data:       &particleData{},
	}
}

func simpleUpdater(data *particleData, elapsed float32, camera components.Viewable) {
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
		data.distanceToCamera[i] = camera.Pos().Sub(data.position[i]).Len()
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

var instanceData = make([]float32, MaxParticles*InstanceDataLength)

func (s *ParticleSystem) Draw(camera components.Viewable, shader components.Shader, state components.RenderState) {

	elapsed := float32(s.timeElapsed)
	toCreate := s.calculateToCreate()

	posX := s.Transform().Pos()[0]
	posY := s.Transform().Pos()[1]
	posZ := s.Transform().Pos()[2]
	for i := 0; i < int(toCreate); i++ {
		colour := [3]float32{2, 2, rand.Float32()*7 + 3}
		s.data.add([3]float32{posX, posY, posZ}, [3]float32{rand.Float32()*1 - 0.5, rand.Float32() * 15, rand.Float32()*1 - 0.5}, colour, rand.Float32()*0.05+0.025, 0, 0.5, rand.Float32()*4+1)
	}

	simpleUpdater(s.data, elapsed, camera)
	sort.Sort(s.data)

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
		instanceData[count+3] = s.data.scale[i]
		copy(instanceData[count+4:], s.data.colour[i][0:3])
		instanceData[count+7] = s.data.transparency[i]
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

func (p *ParticleModel) AABB() components.AABB {
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

func (m *ParticleModel) Update(time.Duration) {
	//fmt.Println("update")
}

func (p *ParticleModel) Unbind() {
	p.mesh.Unbind()
}

type ParticleMesh struct {
	vao uint32
}

func (p *ParticleMesh) AABB() components.AABB {
	return &physics.AABB{}
}

func (p *ParticleMesh) Bind() {
	gl.BindVertexArray(p.vao)
	debug.AddVertexBind()
	gl.EnableVertexAttribArray(0)
	gl.EnableVertexAttribArray(1)
	gl.EnableVertexAttribArray(2)
}

func (p *ParticleMesh) Draw() {
	gl.DrawArraysInstanced(gl.TRIANGLE_STRIP, 0, 4, 1)
	debug.Drawcall()
}

func (p *ParticleMesh) Unbind() {
	gl.DisableVertexAttribArray(2)
	gl.DisableVertexAttribArray(1)
	gl.DisableVertexAttribArray(0)
	gl.BindVertexArray(0)
	debug.AddVertexBind()
}
