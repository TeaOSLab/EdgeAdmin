// Copyright 2022 GoEdge CDN goedge.cdn@gmail.com. All rights reserved. Official site: https://goedge.cn .

package utils_test

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/utils"
	"testing"
	"time"
)

func TestNewUpgradeManager(t *testing.T) {
	var manager = utils.NewUpgradeManager("admin", "")

	var ticker = time.NewTicker(2 * time.Second)
	go func() {
		for range ticker.C {
			if manager.IsDownloading() {
				t.Logf("%.2f%%", manager.Progress()*100)
			}
		}
	}()

	/**go func() {
		time.Sleep(5 * time.Second)
		if manager.IsDownloading() {
			t.Log("cancel downloading")
			_ = manager.Cancel()
		}
	}()**/

	err := manager.Start()
	if err != nil {
		t.Fatal(err)
	}
}
