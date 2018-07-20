package core

import (
	"testing"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/stojg/graphics/lib/components"
	"github.com/stojg/graphics/lib/physics"
	"github.com/stojg/graphics/lib/resources"
)

func TestGameObject_SetModel(t *testing.T) {
	root := NewGameObject()

	first := NewGameObject()
	root.AddChild(first)

	mesh := resources.NewMesh()
	mtrl := resources.NewMaterial()
	model := components.NewModel(mesh, mtrl)

	objA := NewGameObject()
	objA.SetModel(model)
	first.AddChild(objA)

	objB := NewGameObject()
	objB.SetModel(model)
	first.AddChild(objB)

	cam := components.NewCamera(65, 10, 10, 0.1, 100)
	root.RenderAll(cam, &FakeShader{}, &FakeEngine{})

	modelObjects := root.AllModels()
	objects, ok := modelObjects[model]
	if !ok {
		t.Fatalf("Expected to find model in model entity list")
	}

	if len(objects) != 2 {
		t.Fatalf("Expected 2 GameObjects with the same model, got %d", len(objects))
	}

}

type FakeShader struct{}

func (FakeShader) Bind() {
	panic("implement me")
}

func (FakeShader) UpdateUniforms(components.Material, components.RenderState) {
	panic("implement me")
}

func (FakeShader) UpdateTransform(*physics.Transform, components.RenderState) {
	panic("implement me")
}

func (FakeShader) UpdateUniform(name string, value interface{}) {
	panic("implement me")
}

type FakeEngine struct{}

func (FakeEngine) AddCamera(camera components.Viewable) {
	panic("implement me")
}

func (FakeEngine) MainCamera() components.Viewable {
	panic("implement me")
}

func (FakeEngine) SamplerSlot(name string) uint32 {
	panic("implement me")
}

func (FakeEngine) SetSamplerSlot(name string, slot uint32) {
	panic("implement me")
}

func (FakeEngine) AddLight(light components.Light) {
	panic("implement me")
}

func (FakeEngine) Lights() []components.Light {
	panic("implement me")
}

func (FakeEngine) SetActiveLight(light components.Light) {
	panic("implement me")
}

func (FakeEngine) ActiveLight() components.Light {
	panic("implement me")
}

func (FakeEngine) Texture(string) components.Texture {
	panic("implement me")
}

func (FakeEngine) SetTexture(string, components.Texture) {
	panic("implement me")
}

func (FakeEngine) Vector3f(string) mgl32.Vec3 {
	panic("implement me")
}

func (FakeEngine) SetVector3f(string, mgl32.Vec3) {
	panic("implement me")
}

func (FakeEngine) Integer(string) int32 {
	panic("implement me")
}

func (FakeEngine) SetInteger(string, int32) {
	panic("implement me")
}

func (FakeEngine) Float(string) float32 {
	panic("implement me")
}

func (FakeEngine) SetFloat(string, float32) {
	panic("implement me")
}
