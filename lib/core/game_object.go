package core

import (
	"sort"
	"time"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/stojg/geode/lib/components"
	"github.com/stojg/geode/lib/physics"
)

func NewGameObject(rtype int) *GameObject {
	return &GameObject{
		rtype:         rtype,
		transform:     physics.NewTransform(),
		modelEntities: make(map[components.Model][]components.Object),
		childRtypes:   make(map[int][]components.Object),
	}
}

type GameObject struct {
	rtype         int
	children      []components.Object
	childRtypes   map[int][]components.Object
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

func (g *GameObject) Type() int {
	return g.rtype
}

func (g *GameObject) SetModel(model components.Model) {
	g.model = model
}

func (g *GameObject) AddChild(child components.Object) {
	g.children = append(g.children, child)
	g.childRtypes[child.Type()] = append(g.childRtypes[child.Type()], child)
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

func (g *GameObject) AllChildren() []components.Object {
	var a []components.Object
	for _, c := range g.children {
		a = append(a, c)
		a = append(a, c.AllChildren()...)
	}
	return a
}

func (g *GameObject) Update(elapsed time.Duration) {
	g.Transform().Update()
	for _, c := range g.components {
		c.Update(elapsed)
	}
	if g.model != nil {
		g.model.Update(elapsed)
	}
	for _, o := range g.children {
		o.Update(elapsed)
	}
}

type ByModel []components.Object

func (a ByModel) Len() int      { return len(a) }
func (a ByModel) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByModel) Less(i, j int) bool {
	return a[i].Model() != a[j].Model()
}

func (g *GameObject) Render(camera components.Viewable, shader components.Shader, state components.RenderState, rtype int) {
	list := g.AllChildren()
	sort.Sort(ByModel(list))

	shader.Bind()
	var currentModel components.Model
	for i := range list {
		if !list[i].IsType(rtype) {
			continue
		}
		if !list[i].IsVisible(camera) {
			continue
		}
		if currentModel != list[i].Model() {
			list[i].Model().Bind(shader, state)
			currentModel = list[i].Model()
		}
		list[i].Draw(camera, shader, state)
	}
	shader.Unbind()
}

func (g *GameObject) Draw(camera components.Viewable, shader components.Shader, state components.RenderState) {
	shader.UpdateTransform(g.Transform(), state)
	g.Model().Draw()
}

func (g *GameObject) IsVisible(camera components.Viewable) bool {
	if g.model == nil {
		return false
	}

	var bc mgl32.Vec3
	var br mgl32.Vec3

	for i := 0; i < 3; i++ {
		bc[i] = g.transform.Transformation().At(i, 3)
		for j := 0; j < 3; j++ {
			bc[i] += g.transform.Transformation().At(i, j) * g.model.AABB().C()[j]
			br[i] += abs(g.transform.Transformation().At(i, j)) * g.model.AABB().R()[i]
		}
	}
	b := &physics.AABB{}
	b.SetC(bc)
	b.SetR(br)
	return IsVisible(camera.Planes(), b, g.transform.Transformation())
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
