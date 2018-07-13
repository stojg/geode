package ecs

import (
	"reflect"
)

type entity32 struct {
	components     []Component
	componentTypes uint32
}

func (e *entity32) indexOfComponent(componentTyp uint32) int {
	//position value is 1, 2, 3, ...
	//in order to conver it into index, we need to decrement it by one
	position := calcBitIndex(e.componentTypes, componentTyp)
	return int(position - 1)
}

func (e *entity32) Component(component Component) Component {
	//make sure the component exists in list

	componentType := componentTypes[reflect.TypeOf(component)]

	if e.componentTypes&componentType == 0 {
		return nil
	}

	index := e.indexOfComponent(componentType)

	return e.components[index]
}

func (e *entity32) AddComponent(component Component) Entity {
	//make sure that component always passes as non nil value
	if component == nil {
		return e
	}

	t := reflect.TypeOf(component)
	cid, ok := componentTypes[t]
	if !ok {
		cid = nextComponentType
		componentTypes[t] = cid
		nextComponentType++
	}
	component.SetComponentType(cid)
	componentType := component.ComponentType()

	//component already inside this entity32
	if e.componentTypes&componentType != 0 {
		return e
	}

	e.componentTypes = e.componentTypes | componentType

	index := e.indexOfComponent(componentType)

	//insert the new component into right index
	e.components = e.components[0 : len(e.components)+1]
	copy(e.components[index:], e.components[index+1:])
	e.components[index] = component
	return e
}

func (e *entity32) RemoveComponent(componentType uint32) {
	//component doesn't have that component
	if e.componentTypes&componentType == 0 {
		return
	}

	index := e.indexOfComponent(componentType)

	//deleting the component from list
	copy(e.components[index:], e.components[index+1:])
	e.components[len(e.components)-1] = nil
	e.components = e.components[:len(e.components)-1]
}

func (e *entity32) HasComponentTypes(componentTypes uint32) bool {
	return (e.componentTypes & componentTypes) == componentTypes
}

func NewEntity32() Entity {
	return &entity32{
		components: make([]Component, 0, 16),
	}
}
