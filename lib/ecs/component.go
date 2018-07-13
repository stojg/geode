package ecs

import "reflect"

var nextComponentType uint32 = 1
var componentTypes = make(map[reflect.Type]uint32, 0)

type BaseComponent struct {
	cid uint32
}

func (p *BaseComponent) ComponentType() uint32 {
	return p.cid
}

func (p *BaseComponent) SetComponentType(cid uint32) {
	p.cid = cid
}
