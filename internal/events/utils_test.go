package events

import "testing"

func TestOn(t *testing.T) {
	On("hello", func() {
		t.Log("world")
	})
	On("hello", func() {
		t.Log("world2")
	})
	On("hello2", func() {
		t.Log("world2")
	})
	Notify("hello")
}
