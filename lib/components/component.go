package components

type UniformUpdater interface {
}

type Drawable interface {
	Draw()
}

type Bindable interface {
	Bind()
}

type Component interface {
	Update(float32)
	Input(float32)
	Render(Bindable, UniformUpdater)
}

type BaseComponent struct {
	//Transform
}

func (m *BaseComponent) Render(Bindable, UniformUpdater) {}
func (m *BaseComponent) Input(float32)  {}
func (m *BaseComponent) Update(float32) {}

//type Transform struct{}
