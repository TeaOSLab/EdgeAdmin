// Copyright 2021 GoEdge CDN goedge.cdn@gmail.com. All rights reserved.

package goman

import (
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	New(func() {
		t.Log("Hello")

		t.Log(List())
	})

	time.Sleep(1 * time.Second)
	t.Log(List())

	time.Sleep(1 * time.Second)
}

func TestNewWithArgs(t *testing.T) {
	NewWithArgs(func(args ...interface{}) {
		t.Log(args[0], args[1])
	}, 1, 2)
	time.Sleep(1 * time.Second)
}
