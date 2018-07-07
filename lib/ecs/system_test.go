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

// BenchmarkAddSystem-8   	    3000	    521159 ns/op	      25 B/op	       2 allocs/op
// BenchmarkAddSystem-8   	    5000	    291462 ns/op	   80498 B/op	    2005 allocs/op
// BenchmarkAddSystem-8   	   10000	    188817 ns/op	   32899 B/op	      17 allocs/op
// BenchmarkAddSystem-8   	   10000	    105271 ns/op	   32897 B/op	      17 allocs/op
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

func MoveSystem(elapsed float64, pos []*Pos, spd []*Speed) {
	for i := range pos {
		pos[i].X += spd[i].X
	}
	//pos.X += spd.X
}
