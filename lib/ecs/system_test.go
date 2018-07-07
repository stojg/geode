package ecs

import (
	"testing"
)

type Pos struct {
	BaseComponent
	X int
}
type Speed struct {
	BaseComponent
	X int
}

func TestSystem(t *testing.T) {
	Reset()

	a1 := NewEntity()
	a1Pos := &Pos{}
	a1.Add(a1Pos)
	a1.Add(&Speed{X: 1})

	a2 := NewEntity()
	a2Pos := &Pos{}
	a2.Add(a2Pos)
	a2.Add(&Speed{X: -1})

	a3 := NewEntity()
	a3.Add(&Pos{})

	AddSystem(MoveSystem, &Pos{}, &Speed{})
	Update(0.5)

	if a1Pos.X != 1 {
		t.Errorf("Expected a1Pos to be 1, got %d", a1Pos.X)
	}

	if a2Pos.X != -1 {
		t.Errorf("Expected a1Pos to be -1, got %d", a2Pos.X)
	}
}

func BenchmarkAddSystem(b *testing.B) {
	Reset()

	for i := 0; i < 1000; i++ {
		a1 := NewEntity()
		a1.Add(&Pos{})
		a1.Add(&Speed{X: 1})
	}

	AddSystem(MoveSystem, &Pos{}, &Speed{})

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Update(0.5)
	}
}

func MoveSystem(elapsed float64, pos *Pos, spd *Speed) {
	pos.X += spd.X

}
