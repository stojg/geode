package rendering

import (
	"fmt"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/stojg/geode/lib/buffers"
	"github.com/stojg/geode/lib/components"
)

func NewRenderState() *RenderState {
	s := &RenderState{
		samplerMap:    make(map[string]uint32),
		textures:      make(map[string]components.Texture),
		uniforms3f:    make(map[string]mgl32.Vec3),
		uniformsI:     make(map[string]int32),
		uniformsFloat: make(map[string]float32),
	}

	size := 3 * buffers.SizeOfMat4
	s.uboMatrices = buffers.CreateEmptyUBO(size)
	// @todo, this is hacky since the index must be the same as in shader binding, so really, we need to query it per
	// @todo shader and manage it per shader, but we count on there only being one UBO block for now
	index := uint32(0)
	gl.BindBufferRange(gl.UNIFORM_BUFFER, index, s.uboMatrices, 0, size)

	return s
}

type RenderState struct {
	mainCamera    components.Viewer
	uboMatrices   uint32
	lights        []components.Light
	activeLight   components.Light
	samplerMap    map[string]uint32
	textures      map[string]components.Texture
	uniforms3f    map[string]mgl32.Vec3
	uniformsI     map[string]int32
	uniformsFloat map[string]float32
}

func (e *RenderState) Update() {
	view := e.Camera().View()
	invView := view.Inv()
	projection := e.Camera().Projection()
	gl.BindBuffer(gl.UNIFORM_BUFFER, e.uboMatrices)
	gl.BufferSubData(gl.UNIFORM_BUFFER, 0, buffers.SizeOfMat4, gl.Ptr(&view[0]))
	gl.BufferSubData(gl.UNIFORM_BUFFER, buffers.SizeOfMat4, buffers.SizeOfMat4, gl.Ptr(&invView[0]))
	gl.BufferSubData(gl.UNIFORM_BUFFER, 2*buffers.SizeOfMat4, buffers.SizeOfMat4, gl.Ptr(&projection[0]))
	gl.BindBuffer(gl.UNIFORM_BUFFER, 0)

	// @todo add this to the ubo
	e.SetVector3f("x_camPos", e.Camera().Pos())
}

func (e *RenderState) SetActiveLight(light components.Light) {
	e.activeLight = light
}

func (e *RenderState) ActiveLight() components.Light {
	return e.activeLight
}

func (e *RenderState) Lights() []components.Light {
	return e.lights
}

func (e *RenderState) AddLight(l components.Light) {
	e.lights = append(e.lights, l)
}

func (e *RenderState) SetTexture(name string, texture components.Texture) {
	e.textures[name] = texture
}

func (e *RenderState) Texture(name string) components.Texture {
	v, ok := e.textures[name]
	if !ok {
		panic(fmt.Sprintf("Texture, Could not find texture '%s'\n", name))
	}
	return v
}

func (e *RenderState) SetInteger(name string, v int32) {
	e.uniformsI[name] = v
}

func (e *RenderState) Integer(name string) int32 {
	v, ok := e.uniformsI[name]
	if !ok {
		panic(fmt.Sprintf("Integer, no value found for uniform '%s'", name))
	}
	return v
}

func (e *RenderState) SetFloat(name string, v float32) {
	e.uniformsFloat[name] = v
}

func (e *RenderState) Float(name string) float32 {
	v, ok := e.uniformsFloat[name]
	if !ok {
		panic(fmt.Sprintf("Float, no value found for uniform '%s'", name))
	}
	return v
}

func (e *RenderState) SetVector3f(name string, v mgl32.Vec3) {
	e.uniforms3f[name] = v
}

func (e *RenderState) Vector3f(name string) mgl32.Vec3 {
	// @todo set value, regardless, this might be an array that isn't used
	v, ok := e.uniforms3f[name]
	if !ok {
		fmt.Printf("Vector3f, no value found for uniform '%s'\n", name)
	}
	return v
}

func (e *RenderState) SetCamera(c components.Viewer) {
	e.mainCamera = c
}

func (e *RenderState) Camera() components.Viewer {
	return e.mainCamera
}

func (e *RenderState) SamplerSlot(samplerName string) uint32 {
	slot, exists := e.samplerMap[samplerName]
	if !exists {
		fmt.Printf("rendering.Renderer tried finding texture slot for %s, failed\n", samplerName)
	}
	return slot
}

func (e *RenderState) AddSamplerSlot(samplerName string) {
	if _, ok := e.samplerMap[samplerName]; !ok {
		e.samplerMap[samplerName] = uint32(len(e.samplerMap))
	}
}
