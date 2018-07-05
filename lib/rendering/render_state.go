package rendering

import (
	"fmt"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/stojg/graphics/lib/components"
)

func NewRenderState() *RenderState {
	samplerMap := make(map[string]uint32)
	samplerMap["albedo"] = 0
	samplerMap["metallic"] = 1
	samplerMap["roughness"] = 2
	samplerMap["normal"] = 3
	samplerMap["x_shadowMap"] = 9
	samplerMap["x_filterTexture"] = 10
	samplerMap["x_filterTexture2"] = 11
	samplerMap["x_filterTexture3"] = 12
	samplerMap["x_filterTexture4"] = 13

	return &RenderState{
		samplerMap:    samplerMap,
		textures:      make(map[string]components.Texture),
		uniforms3f:    make(map[string]mgl32.Vec3),
		uniformsI:     make(map[string]int32),
		uniformsFloat: make(map[string]float32),
	}
}

type RenderState struct {
	mainCamera    components.Viewable
	lights        []components.Light
	activeLight   components.Light
	samplerMap    map[string]uint32
	textures      map[string]components.Texture
	uniforms3f    map[string]mgl32.Vec3
	uniformsI     map[string]int32
	uniformsFloat map[string]float32
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

func (e *RenderState) AddCamera(c components.Viewable) {
	e.mainCamera = c
}

func (e *RenderState) MainCamera() components.Viewable {
	return e.mainCamera
}

func (e *RenderState) SamplerSlot(samplerName string) uint32 {
	slot, exists := e.samplerMap[samplerName]
	if !exists {
		fmt.Printf("rendering.Engine tried finding texture slot for %s, failed\n", samplerName)
	}
	return slot
}

func (e *RenderState) SetSamplerSlot(samplerName string, slot uint32) {
	e.samplerMap[samplerName] = slot
}
