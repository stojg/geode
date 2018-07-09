package ecs

import (
	"reflect"
)

const maxFunctionArguments = 8

// caches
var (
	in               [maxFunctionArguments]reflect.Value
	componentList    [maxFunctionArguments]reflect.Value
	entityComponents [maxFunctionArguments]int
)

func (e *ECS) AddSystem(system interface{}, components ...interface{}) {
	method := reflect.ValueOf(system)

	id := method.Pointer()
	e.methodPointerToMethod[id] = method
	e.methods = append(e.methods, id)

	for _, comp := range components {
		tid := e.addComponentType(comp.(Component))
		e.methodComponents[id] = append(e.methodComponents[id], tid)
	}

	num := method.Type().NumIn()

	e.systemToIn[id] = make([]reflect.Type, num)
	for i := 1; i < num; i++ {
		e.systemToIn[id][i] = method.Type().In(i)
	}
}

func (e *ECS) Update(elapsed float64) {

	// the first argument to an update function should always be elapsed
	in[0] = reflect.ValueOf(elapsed)

	entities := make([][]int, len(e.allEntityComponents))
	for _, methodId := range e.methods {

		numEntities := 0
		for entityID := range e.allEntityComponents {
			count := 0
			for i, typeID := range e.methodComponents[methodId] {
				for j := range e.allEntityComponentTypes[entityID] {
					if typeID == e.allEntityComponents[entityID][j].TID() {
						entityComponents[i] = e.allEntityComponents[entityID][j].CID()
						count++
						break
					}
				}
			}

			if count != len(e.methodComponents[methodId]) {
				continue
			}

			entities[numEntities] = make([]int, len(entityComponents))
			for i := 0; i < len(entityComponents); i++ {
				entities[numEntities][i] = entityComponents[i]
			}
			numEntities++
		}

		args := e.systemToIn[methodId]

		for i := 1; i < len(args); i++ {
			componentList[i-1] = reflect.MakeSlice(args[i], numEntities, numEntities)
		}

		count := 0
		for _, components := range entities[:numEntities] {
			for i, componentID := range components[:len(args)-1] {
				componentList[i].Index(count).Set(reflect.ValueOf(e.allComponents[componentID]))
			}
			count++
		}
		for i := 1; i < len(args); i++ {
			v, ok := e.allComponentTypes[args[i].Elem()]
			if !ok {
				panic("oh crappers")
			}
			in[i] = componentList[v]
		}
		e.methodPointerToMethod[methodId].Call(in[:len(args)])
	}
}
