package ecs

import "reflect"

// Components are only raw data, ie component Position.x, Position.y ie a struct nothing else/more
// Entities is a collection of Components, Position, Motion, Input nothing else/more
// System, takes a list of Components, ie. All the Position and Motion component in the world, nothing more/else

func New() *ECS {
	return &ECS{
		methods:                 make([]uintptr, 0),
		methodPointerToMethod:   make(map[uintptr]reflect.Value),
		methodComponents:        make(map[uintptr][]int),
		systemToIn:              make(map[uintptr][]reflect.Type),
		allEntityComponents:     make([][]Component, 0),
		allEntityComponentTypes: make([][]int, 0),
		allComponentTypes:       make(map[reflect.Type]int, 0),
		allComponents:           make(map[int]Component),
	}
}

type ECS struct {
	methods                 []uintptr
	methodPointerToMethod   map[uintptr]reflect.Value
	methodComponents        map[uintptr][]int
	systemToIn              map[uintptr][]reflect.Type
	nextEntityID            Entity
	allEntityComponents     [][]Component
	allEntityComponentTypes [][]int
	allComponentTypes       map[reflect.Type]int
	allComponents           map[int]Component
}
