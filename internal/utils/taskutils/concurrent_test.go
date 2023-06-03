// Copyright 2023 GoEdge CDN goedge.cdn@gmail.com. All rights reserved. Official site: https://goedge.cn .

package taskutils_test

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/utils/taskutils"
	"testing"
)

func TestRunConcurrent(t *testing.T) {
	err := taskutils.RunConcurrent([]string{"a", "b", "c", "d", "e"}, 3, func(task any) {
		t.Log("run", task)
	})
	if err != nil {
		t.Fatal(err)
	}
}
