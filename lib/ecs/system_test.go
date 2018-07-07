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
	e := New()
	a1 := e.NewEntity()
	a1Pos := &Pos{}
	e.Add(a1, a1Pos)
	e.Add(a1, &Speed{X: 1})

	a2 := e.NewEntity()
	a2Pos := &Pos{}
	e.Add(a2, a2Pos)
	e.Add(a2, &Speed{X: -1})

	a3 := e.NewEntity()
	e.Add(a3, &Pos{})

	e.AddSystem(MoveSystem, &Pos{}, &Speed{})
	e.AddSystem(RenderSystem, &Pos{})
	e.Update(0.5)

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
// BenchmarkAddSystem-8   	   10000	    105271 ns/op	   32897 B/op	      17 allocs/op (buggy)
// BenchmarkAddSystem-8   	    2000	    636289 ns/op	  185624 B/op	    6031 allocs/op
// BenchmarkAddSystem-8   	    2000	    696510 ns/op	  422113 B/op	    3070 allocs/op
// BenchmarkAddSystem-8   	    3000	    426758 ns/op	  187206 B/op	    3031 allocs/op
// BenchmarkAddSystem-8   	    5000	    303953 ns/op	  171189 B/op	    2031 allocs/op
func BenchmarkAddSystem(b *testing.B) {
	e := New()

	for i := 0; i < 1000; i++ {
		a1 := e.NewEntity()
		e.Add(a1, &Pos{})
		e.Add(a1, &Speed{X: 1})
	}

	e.AddSystem(MoveSystem, &Pos{}, &Speed{})
	e.AddSystem(RenderSystem, &Pos{})

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		e.Update(0.5)
	}
}

func MoveSystem(elapsed float64, pos []*Pos, spd []*Speed) {
	for i := range pos {
		pos[i].X += spd[i].X
	}
}

func RenderSystem(elapsed float64, pos []*Pos) {
	//for i := range pos {
	//fmt.Println(pos[i].X)
	//}
}
