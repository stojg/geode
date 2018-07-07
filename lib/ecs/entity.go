package ecs

import "fmt"

type Entity int

func (e Entity) Add(c Component) {
	addComponent(c)
	// check that we don'existing add multiple components of the same type
	for _, existing := range allEntityComponents[e] {
		if existing.TID() == c.TID() {
			panic(fmt.Sprintf("Entity %d already have component of type %d", e, existing))
		}
	}
	allEntityComponentTypes[e] = append(allEntityComponentTypes[e], c.TID())
	allEntityComponents[e] = append(allEntityComponents[e], c)
	c.setEID(e)
}

func (e Entity) getComponentTypes() []int {
	return allEntityComponentTypes[e]
}

func (e Entity) getComponent(compType int) Component {
	v, ok := allEntityComponents[e]
	if !ok {
		panic(fmt.Sprintf("Entity %d cant be found in entityCommponents list", e))
	}

	for _, c := range v {
		if c.TID() == compType {
			return c
		}
	}
	panic(fmt.Sprintf("Component type %d cant be found in entityCommponents list", compType))
	return nil
}

func (e Entity) getComponents() []Component {
	v, ok := allEntityComponents[e]
	if !ok {
		panic(fmt.Sprintf("Entity %d cant be found in entityCommponents list", e))
	}
	return v
}

func NewEntity() Entity {
	eid := nextEntityID
	nextEntityID++
	allEntityComponents[eid] = make([]Component, 0)
	return eid
}
