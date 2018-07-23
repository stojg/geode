package core

import (
	"time"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/stojg/graphics/lib/components"
	"github.com/stojg/graphics/lib/physics"
)

func NewGameObject(rtype int) *GameObject {
	return &GameObject{
		rtype:         rtype,
		transform:     physics.NewTransform(),
		modelEntities: make(map[components.Model][]components.Object),
	}
}

type GameObject struct {
	rtype         int
	children      []components.Object
	components    []components.Component
	transform     *physics.Transform
	state         components.RenderState
	model         components.Model
	modelEntities map[components.Model][]components.Object
}

func (g *GameObject) Model() components.Model {
	return g.model
}

func (g *GameObject) IsType(a int) bool {
	return g.rtype&a != 0
}

func (g *GameObject) SetModel(model components.Model) {
	g.model = model
}

func (g *GameObject) AddChild(child components.Object) {
	g.children = append(g.children, child)
	child.SetState(g.state)
	child.Transform().SetParent(g.Transform())
	if child.Model() != nil {
		m := child.Model()
		g.modelEntities[m] = append(g.modelEntities[m], child)
	}
}

func (g *GameObject) AddComponent(component components.Component) {
	component.SetParent(g)
	g.components = append(g.components, component)
}

func (g *GameObject) Input(elapsed time.Duration) {
	for _, c := range g.components {
		c.Input(elapsed)
	}
	for _, o := range g.children {
		o.Input(elapsed)
	}
}

func (g *GameObject) AllModels() map[components.Model][]components.Object {
	a := make(map[components.Model][]components.Object)
	for _, c := range g.children {
		for k, v := range c.AllModels() {
			a[k] = append(a[k], v...)
		}
	}
	for k, v := range g.modelEntities {
		a[k] = append(a[k], v...)
	}
	return a
}

func (g *GameObject) Update(elapsed time.Duration) {
	g.Transform().Update()
	for _, c := range g.components {
		c.Update(elapsed)
	}
	for _, o := range g.children {
		o.Update(elapsed)
	}
}

func (g *GameObject) Render(camera components.Viewable, shader components.Shader, state components.RenderState, rtype int) {
	list := g.AllModels()

	for model, objects := range list {
		var visible []int
		for i := range objects {

			if !objects[i].IsType(rtype) {
				continue
			}
			if objects[i].IsVisible(camera) {
				visible = append(visible, i)
			}
		}

		if len(visible) == 0 {
			continue
		}

		shader.Bind()
		model.Bind(shader, state)
		for _, i := range visible {
			objects[i].Draw(camera, shader, state)
		}
		model.Unbind()
	}
}

func (g *GameObject) Draw(camera components.Viewable, shader components.Shader, state components.RenderState) {
	shader.UpdateTransform(g.Transform(), state)
	g.Model().Draw()
}

func (g *GameObject) render(shader components.Shader, state components.RenderState) {
	shader.UpdateTransform(g.Transform(), state)
	g.model.Draw()
}

func (g *GameObject) IsVisible(camera components.Viewable) bool {
	return IsVisible(camera.Planes(), g.model.AABB(), g.transform.Transformation())
}

func (g *GameObject) Transform() *physics.Transform {
	return g.transform
}

func (g *GameObject) SetPos(x, y, z float32) {
	g.transform.SetPos(mgl32.Vec3{x, y, z})
}

func (g *GameObject) SetScale(x, y, z float32) {
	g.transform.SetScale(mgl32.Vec3{x, y, z})
}

func (g *GameObject) Rotate(axis mgl32.Vec3, radians float32) {
	g.transform.Rotate(axis, radians)
}

func (g *GameObject) SetState(state components.RenderState) {
	if g.state != state {
		g.state = state
		for _, c := range g.components {
			c.AddToEngine(state)
		}
		for _, c := range g.children {
			c.SetState(state)
		}
	}
}
