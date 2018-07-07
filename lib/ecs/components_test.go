package ecs_test

import (
	"testing"

	"github.com/stojg/graphics/lib/ecs"
)

func TestNewComponent(t *testing.T) {
	e := ecs.New()
	entity := e.NewEntity()

	type Position struct {
		ecs.BaseComponent
		x, y, z float32
	}

	pos := &Position{}
	e.Add(entity, pos)
	if pos.TID() != 0 {
		t.Errorf("pos TID should be 0, got %d", pos.TID())
	}
	if pos.CID() != 0 {
		t.Errorf("pos CID should be 0, got %d", pos.CID())
	}

	type Movement struct {
		ecs.BaseComponent
		x, y, z float32
	}

	mv := &Movement{x: 0}
	e.Add(entity, mv)
	if mv.TID() != 1 {
		t.Errorf("mv TID should be 1, got %d", mv.TID())
	}
	if mv.TID() != 1 {
		t.Errorf("mv CID should be 1, got %d", mv.CID())
	}

	otherEntity := e.NewEntity()

	mv2 := &Movement{x: 0}
	e.Add(otherEntity, mv2)
	if mv2.TID() != 1 {
		t.Errorf("mv2 TID should be 1, got %d", mv2.TID())
	}
	if mv2.CID() != 2 {
		t.Errorf("mv2 CID should be 2, got %d", mv2.CID())
	}
}
