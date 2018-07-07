package ecs

import "fmt"

type Entity int

func (e Entity) Add(c Component) {
	addComponent(c)
	for _, existing := range allEntityComponents[e] {
		if existing.TID() == c.TID() {
			panic(fmt.Sprintf("Entity %d already have component of type %d", e, existing))
		}
	}
	if len(allEntityComponentTypes) <= int(e) {
		allEntityComponentTypes = append(allEntityComponentTypes, []int{c.TID()})
	} else {
		allEntityComponentTypes[e] = append(allEntityComponentTypes[e], c.TID())
	}
	allEntityComponents[e] = append(allEntityComponents[e], c)
	c.setEID(e)
}

func (e Entity) getComponentTypes() []int {
	return allEntityComponentTypes[e]
}

func (e Entity) getComponent(compType int) Component {
	return allEntityComponents[e][compType]
}

func (e Entity) getComponents() []Component {
	return allEntityComponents[e]
}

func NewEntity() Entity {
	eid := nextEntityID
	nextEntityID++
	allEntityComponents = append(allEntityComponents, []Component{})
	return eid
}
