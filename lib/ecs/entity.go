package ecs

import "fmt"

type Entity int

func (d *ECS) Add(e Entity, c Component) {
	d.addComponent(c)
	for _, existing := range d.allEntityComponents[e] {
		if existing.TID() == c.TID() {
			panic(fmt.Sprintf("Entity %d already have component of type %d", e, existing))
		}
	}
	if len(d.allEntityComponentTypes) <= int(e) {
		d.allEntityComponentTypes = append(d.allEntityComponentTypes, []int{c.TID()})
	} else {
		d.allEntityComponentTypes[e] = append(d.allEntityComponentTypes[e], c.TID())
	}
	d.allEntityComponents[e] = append(d.allEntityComponents[e], c)
	c.setEID(e)
}

func (e *ECS) NewEntity() Entity {
	eid := e.nextEntityID
	e.nextEntityID++
	e.allEntityComponents = append(e.allEntityComponents, []Component{})
	return eid
}
