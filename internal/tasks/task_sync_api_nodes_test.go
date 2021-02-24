package tasks

import (
	_ "github.com/iwind/TeaGo/bootstrap"
	"testing"
)

func TestSyncAPINodesTask_Loop(t *testing.T) {
	task := NewSyncAPINodesTask()
	err := task.Loop()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("ok")
}
