package core

import (
	"time"

	"github.com/stojg/graphics/lib/components"
	"github.com/stojg/graphics/lib/physics"
)

func NewGameObject() *GameObject {
	return &GameObject{
		transform:     physics.NewTransform(),
		modelEntities: make(map[*components.Model][]*GameObject),
	}
}

type GameObject struct {
	children   []*GameObject
	components []components.Component
	transform  *physics.Transform
	state      components.RenderState

	model         *components.Model
	modelEntities map[*components.Model][]*GameObject
}

func (g *GameObject) Model() *components.Model {
	return g.model
}

func (g *GameObject) SetModel(model *components.Model) {
	g.model = model
}

func (g *GameObject) AddChild(child *GameObject) {
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
}

func (g *GameObject) InputAll(elapsed time.Duration) {
	g.Input(elapsed)
	for _, o := range g.children {
		o.InputAll(elapsed)
	}
}

func (g *GameObject) Update(elapsed time.Duration) {
	g.Transform().Update()
	for _, c := range g.components {
		c.Update(elapsed)
	}
}

func (g *GameObject) AllModels() map[*components.Model][]*GameObject {
	a := make(map[*components.Model][]*GameObject)
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

func (g *GameObject) UpdateAll(elapsed time.Duration) {
	g.Update(elapsed)
	for _, o := range g.children {
		o.UpdateAll(elapsed)
	}
}

func (g *GameObject) Render(shader components.Shader, state components.RenderState) {
	shader.UpdateTransform(g.Transform(), state)
	g.model.Draw()
}

func (g *GameObject) RenderAll(camera components.Viewable, shader components.Shader, state components.RenderState) {
	list := g.AllModels()
	shader.Bind()
	for model, objects := range list {
		model.Bind(shader, state)
		for _, object := range objects {
			if IsVisible(camera.Planes(), object.model.AABB(), object.transform.Transformation()) {
				object.Render(shader, state)
			}
		}
		model.Unbind()
	}
}

func (g *GameObject) Transform() *physics.Transform {
	return g.transform
}

func (g *GameObject) AllTransforms() []*physics.Transform {
	var l []*physics.Transform
	for _, c := range g.children {
		l = append(l, c.Transform())
	}
	return l
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
