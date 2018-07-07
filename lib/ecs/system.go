package ecs

import (
	"reflect"
)

func AddSystem(system interface{}, components ...Component) {
	method := reflect.ValueOf(system)
	for _, comp := range components {
		tid := addComponentType(comp)
		systemComponents[method] = append(systemComponents[method], tid)
	}

	in := make([]reflect.Type, method.Type().NumIn())
	for i := 0; i < method.Type().NumIn(); i++ {
		in[i] = method.Type().In(i)
	}
	systemToIn[method] = in
}

func Update(elapsed float64) {
	in := make([]reflect.Value, 16)
	in[0] = reflect.ValueOf(elapsed)
	listB := make([]reflect.Value, 16)

	for method, args := range systemToIn {
		for i := 1; i < len(args); i++ {
			v, _ := allComponentTypes[args[i].Elem()]
			listB[v] = reflect.MakeSlice(args[i], 0, len(getAllEntities()))
		}
		for entity, components := range getAllEntities() {
			if !canEntityBeUpdated(entity, systemComponents[method]) {
				continue
			}
			for _, component := range components {
				v := reflect.ValueOf(component)
				listB[component.TID()] = reflect.Append(listB[component.TID()], v)
			}
		}
		for i := 1; i < len(args); i++ {
			v, _ := allComponentTypes[args[i].Elem()]
			in[i] = listB[v]
		}
		method.Call(in[:len(args)])
	}
}

func canEntityBeUpdated(entity Entity, componentTypes []int) bool {
	count := 0
	for _, typeID := range componentTypes {
		for _, entityComponentID := range entity.getComponentTypes() {
			if typeID == entityComponentID {
				count++
				break
			}
		}
	}
	return count == len(componentTypes)
}
