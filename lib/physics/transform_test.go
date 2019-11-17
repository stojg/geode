package physics

import (
	"testing"

	"github.com/go-gl/mathgl/mgl32"
)

func TestTransform_Update(t *testing.T) {
	tr := NewTransform()

	tr.Update()

	if tr.dirty {
		t.Error("Expected dirty to be false")
	}

	tr.SetPos(mgl32.Vec3{1, 0, 0})

	if !tr.dirty {
		t.Error("Expected dirty to be true")
	}

	tr.Update()
	if tr.dirty {
		t.Error("Expected dirty to be false after update")
	}
}
