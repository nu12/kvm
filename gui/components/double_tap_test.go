package components

import (
	"testing"

	"fyne.io/fyne/v2/test"
)

func TestDoubleTap(t *testing.T) {
	c := NewDoubleTapCell()
	if c == nil {
		t.Fatalf("NewDoubleTapCell() returned nil")
	}

}

func TestDoubleTapped(t *testing.T) {
	called := false

	c := NewDoubleTapCell()
	c.DoubleTapFunc = func() {
		called = true
	}

	test.DoubleTap(c)

	if !called {
		t.Error("DoubleTapFunc was not called after DoubleTap")
	}
}
