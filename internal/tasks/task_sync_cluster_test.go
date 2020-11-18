package tasks

import "testing"

func TestSyncClusterTask_loop(t *testing.T) {
	task := NewSyncClusterTask()
	err := task.loop()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("ok")
}
