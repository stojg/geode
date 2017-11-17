package core

import (
	"github.com/stojg/graphics/lib/components"
	"github.com/stojg/graphics/lib/rendering"
)

func NewGameObject() *GameObject {
	return &GameObject{}
}

type GameObject struct {
	children   []*GameObject
	components []components.Component
}

func (g *GameObject) AddChild(child *GameObject) {
	g.children = append(g.children, child)
}

func (g *GameObject) AddComponent(component components.Component) {
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

func (g *GameObject) RenderAll(shader *rendering.Shader, renderingEngine components.UniformUpdater) {
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

func (g *GameObject) Render(shader *rendering.Shader, renderingEngine components.UniformUpdater) {
	for _, c := range g.components {
		c.Render(shader, renderingEngine)
	}
}
