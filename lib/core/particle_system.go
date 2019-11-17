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
	"github.com/stojg/geode/lib/geometry"
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

func (pd *particleData) add(pos, vel, colour [3]float32, scale, rotAngle, gravity, life float32) {
	if pd.aliveCount >= MaxParticles {
		return
	}
	pd.alive[pd.aliveCount] = true
	pd.transparency[pd.aliveCount] = 1.0
	pd.position[pd.aliveCount] = pos
	pd.velocity[pd.aliveCount] = vel
	pd.scale[pd.aliveCount] = scale
	pd.rotation[pd.aliveCount] = rotAngle
	pd.gravity[pd.aliveCount] = gravity * Gravity
	pd.lifeLength[pd.aliveCount] = life
	pd.elapsedTime[pd.aliveCount] = 0
	pd.distanceToCamera[pd.aliveCount] = 0
	pd.colour[pd.aliveCount] = colour
	pd.aliveCount++
}

func (pd *particleData) remove(id int) {
	pd.alive[id] = false
	pd.swap(id, pd.aliveCount-1)
	pd.aliveCount--
}

func (pd *particleData) swap(a, b int) {
	pd.velocity[a], pd.velocity[b] = pd.velocity[b], pd.velocity[a]
	pd.gravity[a], pd.gravity[b] = pd.gravity[b], pd.gravity[a]
	pd.position[a], pd.position[b] = pd.position[b], pd.position[a]
	pd.transparency[a], pd.transparency[b] = pd.transparency[b], pd.transparency[a]
	pd.rotation[a], pd.rotation[b] = pd.rotation[b], pd.rotation[a]
	pd.scale[a], pd.scale[b] = pd.scale[b], pd.scale[a]
	pd.elapsedTime[a], pd.elapsedTime[b] = pd.elapsedTime[b], pd.elapsedTime[a]
	pd.alive[a], pd.alive[b] = pd.alive[b], pd.alive[a]
	pd.lifeLength[a], pd.lifeLength[b] = pd.lifeLength[b], pd.lifeLength[a]
	pd.distanceToCamera[a], pd.distanceToCamera[b] = pd.distanceToCamera[b], pd.distanceToCamera[a]
	pd.colour[a], pd.colour[b] = pd.colour[b], pd.colour[a]
}

func (pd *particleData) Len() int           { return pd.aliveCount }
func (pd *particleData) Swap(i, j int)      { pd.swap(i, j) }
func (pd *particleData) Less(i, j int) bool { return pd.distanceToCamera[i] > pd.distanceToCamera[j] }

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

func simpleUpdater(data *particleData, elapsed float32, camera components.Viewer) {
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

func (s *ParticleSystem) Draw(camera components.Viewer, shader components.Shader, state components.RenderState) {
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

func (s *ParticleSystem) IsVisible(camera components.Viewer) bool {
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

func (pm *ParticleModel) AABB() components.AABB {
	return pm.mesh.AABB()
}

func (pm *ParticleModel) Bind(shader components.Shader, state components.RenderState) {
	shader.UpdateUniforms(pm.material, state)
	pm.mesh.Bind()
}

func (pm *ParticleModel) Material() components.Material {
	return pm.material
}

func (pm *ParticleModel) Draw() {
	pm.mesh.Draw()
}

func (pm *ParticleModel) Update(time.Duration) {
	//fmt.Println("update")
}

func (pm *ParticleModel) Unbind() {
	pm.mesh.Unbind()
}

type ParticleMesh struct {
	vao uint32
}

func (p *ParticleMesh) AABB() components.AABB {
	return &geometry.AABB{}
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
