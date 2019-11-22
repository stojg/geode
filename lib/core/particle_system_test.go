package core

import (
	"testing"
	"time"

	"github.com/stojg/geode/lib/components"
)

func Test_simpleUpdater(t *testing.T) {
	data := &particleData{}

	data.add([3]float32{1, 0, 0}, [3]float32{0, 0, 0}, [3]float32{1, 1, 1}, 0, 1, 1, 1)

	cameraObject := NewGameObject(components.ResourceNA)
	cam := components.NewCamera(75, 320, 240, 0.1, 512)
	cameraObject.AddComponent(cam)

	simpleUpdater(data, 0.1)

	if data.Len() != 1 {
		t.Errorf("expected 1 alive particle")
	}

	if data.elapsedTime[0] != 0.1 {
		t.Errorf("expected particle idx 0 to be alive")
	}

	data.add([3]float32{2, 0, 0}, [3]float32{0, 0, 0}, [3]float32{1, 1, 1}, 0, 1, 1, 1)

	if data.Len() != 2 {
		t.Errorf("expected 2 alive particles")
	}

	if data.elapsedTime[0] != 0.1 {
		t.Errorf("expected particle idx 1 to be alive")
	}

	// now we zoom past the alive state for th first particle
	simpleUpdater(data, 1.0)

	if data.Len() != 1 {
		t.Errorf("expected 1 alive particlesx")
	}

	if data.elapsedTime[0] != 1 {
		t.Errorf("expected particle idx 0 to be %f old, got %f", 1.0, data.elapsedTime[0])
	}

	if data.elapsedTime[1] != 1.1 {
		t.Errorf("expected particle idx 1 to be %f old, got %f", 1.1, data.elapsedTime[1])
	}

	if data.position[0][0] != 2.0 {
		t.Errorf("expected idx 1 to have moved to idx 0 after idx 0 been killed")
	}

	simpleUpdater(data, 0.01)

	if data.elapsedTime[0] != 1.01 {
		t.Errorf("expected particle idx 0 to be %f old, got %f", 1.01, data.elapsedTime[0])
	}

	if data.elapsedTime[1] != 1.1 {
		t.Errorf("expected particle idx 1 to be %f old, got %f", 1.1, data.elapsedTime[1])
	}
}

func Benchmark_simpleUpdater(b *testing.B) {
	data := &particleData{}

	cameraObject := NewGameObject(components.ResourceNA)
	cam := components.NewCamera(75, 320, 240, 0.1, 512)
	cameraObject.AddComponent(cam)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		data.add([3]float32{float32(i), 0, 0}, [3]float32{1, 0, 0}, [3]float32{1, 1, 1}, 0, 1, 1, 50)
		simpleUpdater(data, 1)
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
