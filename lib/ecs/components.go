package ecs

import (
	"reflect"
)

var allComponentTypes = make(map[reflect.Type]int, 0)
var allComponents = make(map[int]Component)
var nextEntityID Entity
var allEntityComponents = make(map[Entity][]Component)
var allEntityComponentTypes = make(map[Entity][]int)
var systemComponents = make(map[reflect.Value][]int)
var systemToIn = make(map[reflect.Value][]reflect.Type)

func getAllEntities() map[Entity][]Component {
	return allEntityComponents
}

func Reset() {
	allComponentTypes = make(map[reflect.Type]int, 0)
	allComponents = make(map[int]Component)
	nextEntityID = 0
	allEntityComponents = make(map[Entity][]Component)
	allEntityComponentTypes = make(map[Entity][]int)
	systemComponents = make(map[reflect.Value][]int)
	systemToIn = make(map[reflect.Value][]reflect.Type)
}

func addComponent(x Component) {
	tid := addComponentType(x)
	x.setTID(tid)
	cid := len(allComponents)
	allComponents[cid] = x
	x.setCID(cid)
}

func addComponentType(x Component) int {
	t := reflect.TypeOf(x)
	v, ok := allComponentTypes[t]
	if !ok {
		v = len(allComponentTypes)
		allComponentTypes[t] = v
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
