package core

import (
	"github.com/stojg/graphics/lib/components"
	"github.com/stojg/graphics/lib/physics"
	"github.com/stojg/graphics/lib/rendering"
)

func NewGameObject() *GameObject {
	return &GameObject{
		transform: physics.NewTransform(),
	}
}

type GameObject struct {
	children   []*GameObject
	components []components.Component
	transform  *physics.Transform
	engine     *Engine
}

func (g *GameObject) AddChild(child *GameObject) {
	g.children = append(g.children, child)
	child.SetEngine(g.engine)
	child.Transform().SetParent(g.Transform())
}

func (g *GameObject) AddComponent(component components.Component) {
	component.SetParent(g)
	g.components = append(g.components, component)
}

func (g *GameObject) InputAll(elapsed float32) {
	g.Input(elapsed)
	for _, o := range g.children {
		o.InputAll(elapsed)
	}
}

func (g *GameObject) UpdateAll(elapsed float32) {
	g.Update(elapsed)
	for _, o := range g.children {
		o.UpdateAll(elapsed)
	}
}

func (g *GameObject) RenderAll(shader *rendering.Shader, renderingEngine components.RenderingEngine) {
	g.Render(shader, renderingEngine)
	for _, o := range g.children {
		o.RenderAll(shader, renderingEngine)
	}
}

func (g *GameObject) Input(elapsed float32) {
	for _, c := range g.components {
		c.Input(elapsed)
	}
}

func (g *GameObject) Update(elapsed float32) {
	for _, c := range g.components {
		c.Update(elapsed)
	}
}

func (g *GameObject) Render(shader *rendering.Shader, renderingEngine components.RenderingEngine) {
	for _, c := range g.components {
		c.Render(shader, renderingEngine)
	}
}

func (g *GameObject) Transform() *physics.Transform {
	return g.transform
}

func (g *GameObject) SetEngine(engine *Engine) {
	//if g.engine != engine {
	g.engine = engine
	for _, c := range g.components {
		c.AddToEngine(engine)
	}
	for _, c := range g.children {
		c.SetEngine(engine)
	}
	//}
}
