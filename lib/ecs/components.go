package ecs

import (
	"reflect"
)

func (e *ECS) addComponent(x Component) {
	tid := e.addComponentType(x)
	x.setTID(tid)
	cid := len(e.allComponents)
	e.allComponents[cid] = x
	x.setCID(cid)
}

func (e *ECS) addComponentType(x Component) int {
	t := reflect.TypeOf(x)
	v, ok := e.allComponentTypes[t]
	if !ok {
		v = len(e.allComponentTypes)
		e.allComponentTypes[t] = v
	}
	return v
}

type Component interface {
	CID() int
	TID() int
	EID() Entity
	setCID(id int)
	setTID(id int)
	setEID(id Entity)
}

type BaseComponent struct {
	cid int
	tid int
	eid Entity
}

func (b *BaseComponent) CID() int {
	return b.cid
}

func (b *BaseComponent) TID() int {
	return b.tid
}

func (b *BaseComponent) EID() Entity {
	return b.eid
}

func (b *BaseComponent) setCID(id int) {
	b.cid = id
}

func (b *BaseComponent) setTID(typeID int) {
	b.tid = typeID
}

func (b *BaseComponent) setEID(entityID Entity) {
	b.eid = entityID
}
