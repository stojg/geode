package ecs_test

import (
	"testing"

	"github.com/stojg/graphics/lib/ecs"
)

func TestSomething(t *testing.T) {
	w := ecs.NewWorld()

	e := ecs.NewEntity32()
	e.AddComponent(&Pos{X: 1.0}).
		AddComponent(&Speed{}).
		AddComponent(&Talk{}).
		AddComponent(&Run{Length: 2.0})
	w.AddEntity(e)

	e2 := ecs.NewEntity32()
	e2.AddComponent(&Pos{X: 10.0}).
		AddComponent(&Run{Length: 20})
	w.AddEntity(e2)

	e3 := ecs.NewEntity32()
	e3.AddComponent(&Pos{X: 10.0})
	w.AddEntity(e3)

	sys := &TestSystem{}
	w.AddSystem(sys)

	w.Update(1)

	if sys.updated != 2 {
		t.Errorf("Expected %d entities to be updated, got %d", 2, sys.updated)
	}

}

func BenchmarkAdd(b *testing.B) {

	w := ecs.NewWorld()

	for i := 0; i < 1000; i++ {
		e := ecs.NewEntity32().
			AddComponent(&Pos{X: 1.0}).
			AddComponent(&Speed{})
		w.AddEntity(e)
	}
	b.ReportAllocs()
	b.ResetTimer()

	sys := &TestSystem{}
	w.AddSystem(sys)

	for i := 0; i < b.N; i++ {
		w.Update(1)
	}

}

type Pos struct {
	ecs.BaseComponent
	X, Y, Z float32
}

type Speed struct {
	ecs.BaseComponent
}

type Talk struct {
	ecs.BaseComponent
}

type Run struct {
	ecs.BaseComponent
	Length float32
}

type TestSystem struct {
	updated int
}

func (s *TestSystem) Update(delta float32, query ecs.Query) {
	for _, en := range query.Entities(&Pos{}, &Run{}) {
		s.updated++
		p := en.Component(&Pos{}).(*Pos)
		r := en.Component(&Run{}).(*Run)
		p.X += r.Length
	}
}

func (s *TestSystem) SystemType() uint32 {
	return 1
}
