package core

import (
	"testing"
	"time"
)

func Test_simpleUpdater(t *testing.T) {
	data := &particleData{}

	data.add([3]float32{1, 0, 0}, [3]float32{0, 0, 0}, 1, 0, 1, 1)

	simpleUpdater(data, 0.1)

	if data.aliveCount != 1 {
		t.Errorf("expected 1 alive particle")
	}

	if !data.alive[0] {
		t.Errorf("expected particle idx 0 to be alive")
	}

	data.add([3]float32{2, 0, 0}, [3]float32{0, 0, 0}, 1, 0, 1, 1)

	if data.aliveCount != 2 {
		t.Errorf("expected 2 alive particles")
	}

	if !data.alive[0] {
		t.Errorf("expected particle idx 1 to be alive")
	}

	simpleUpdater(data, 1.0)

	if data.aliveCount != 1 {
		t.Errorf("expected 1 alive particles")
	}

	if !data.alive[0] {
		t.Errorf("expected particle idx 0 to be alive")
	}

	if data.alive[1] {
		t.Errorf("expected particle idx 1 to be dead")
	}

	if data.position[0][0] != 2.0 {
		t.Errorf("expected idx 1 to have moved to idx 0 after idx 0 been killed")
	}
}

func Test_calculateToCreate(t *testing.T) {
	s := &ParticleSystem{
		perSecond: 10.0,
	}
	s.Update(time.Second)

	if s.timeElapsed != 1.0 {
		t.Errorf("expected %f, got %f", 1.0, s.timeElapsed)
	}

	toCreate := s.calculateToCreate()

	if toCreate != 10 {
		t.Errorf("expected %f, got %f", 10.0, toCreate)
	}
	if s.reminder != 0 {
		t.Errorf("expected %f, got %f", 0.0, s.reminder)
	}

	s.Update(time.Second / 2)
	toCreate = s.calculateToCreate()
	if toCreate != 5 {
		t.Errorf("expected %f, got %f", 5.0, toCreate)
	}
	if s.reminder != 0 {
		t.Errorf("expected %f, got %f", 0.0, s.reminder)
	}

	s.Update(time.Millisecond * 50)
	toCreate = s.calculateToCreate()
	if toCreate != 0 {
		t.Errorf("expected %f, got %f", 0.0, toCreate)
	}
	if s.reminder != 0.5 {
		t.Errorf("expected %f, got %f", 0.5, s.reminder)
	}

	s.Update(time.Millisecond * 700)
	toCreate = s.calculateToCreate()
	if toCreate != 7 {
		t.Errorf("expected %f, got %f", 7.0, toCreate)
	}
	if s.reminder != 0.5 {
		t.Errorf("expected %f, got %f", 0.5, s.reminder)
	}

	s.Update(time.Millisecond * 50)
	toCreate = s.calculateToCreate()
	if toCreate != 1 {
		t.Errorf("expected %f, got %f", 1.0, toCreate)
	}
	if s.reminder != 0.0 {
		t.Errorf("expected %f, got %f", 0.0, s.reminder)
	}
}
