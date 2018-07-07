package ecs

import "reflect"

// Components are only raw data, ie component Position.x, Position.y ie a struct nothing else/more
// Entities is a collection of Components, Position, Motion, Input nothing else/more
// System, takes a list of Components, ie. All the Position and Motion component in the world, nothing more/else

func New() *ECS {
	return &ECS{
		systemComponents:        make(map[reflect.Value][]int),
		systemToIn:              make(map[reflect.Value][]reflect.Type),
		allEntityComponents:     make([][]Component, 0),
		allEntityComponentTypes: make([][]int, 0),
		allComponentTypes:       make(map[reflect.Type]int, 0),
		allComponents:           make(map[int]Component),
	}
}

type ECS struct {
	systemComponents        map[reflect.Value][]int
	systemToIn              map[reflect.Value][]reflect.Type
	nextEntityID            Entity
	allEntityComponents     [][]Component
	allEntityComponentTypes [][]int
	allComponentTypes       map[reflect.Type]int
	allComponents           map[int]Component
}
