package ecs

import (
	"fmt"
	"reflect"
)

func AddSystem(system interface{}, components ...Component) {
	method := reflect.ValueOf(system)
	for _, comp := range components {
		tid := addComponentType(comp)
		systemComponents[method] = append(systemComponents[method], tid)
	}

	num := method.Type().NumIn()
	systemToIn[method] = make([]reflect.Type, num)
	for i := 1; i < num; i++ {
		systemToIn[method][i] = method.Type().In(i)
	}
}

func Update(elapsed float64) {
	in := make([]reflect.Value, 8)
	in[0] = reflect.ValueOf(elapsed)
	listB := make([]reflect.Value, len(allComponentTypes))

	for method, args := range systemToIn {
		count := 0
		var entities []Entity
		for entity := range getAllEntities() {
			if !canEntityBeUpdated(entity, systemComponents[method]) {
				continue
			}
			entities = append(entities, Entity(entity))
		}

		for i := 1; i < len(args); i++ {
			componentID, ok := allComponentTypes[args[i].Elem()]
			if !ok {
				panic(fmt.Sprintf("Can't find component type for %s", args[i].Elem()))
			}
			listB[componentID] = reflect.MakeSlice(args[i], len(entities), len(entities))
		}

		for _, entity := range entities {
			// @todo, only grab the relevant components
			for _, component := range entity.getComponents() {
				v := reflect.ValueOf(component)
				listB[component.TID()].Index(count).Set(v)
			}
			count++
		}
		for i := 1; i < len(args); i++ {
			v, _ := allComponentTypes[args[i].Elem()]
			in[i] = listB[v]
		}
		method.Call(in[:len(args)])
	}
}

func canEntityBeUpdated(entity int, componentTypes []int) bool {
	count := 0
	e := Entity(entity)
	for _, typeID := range componentTypes {
		for _, entityComponentID := range e.getComponentTypes() {
			if typeID == entityComponentID {
				count++
				break
			}
		}
	}
	return count == len(componentTypes)
}
