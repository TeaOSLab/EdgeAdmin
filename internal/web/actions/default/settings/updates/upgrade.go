// Copyright 2023 GoEdge CDN goedge.cdn@gmail.com. All rights reserved. Official site: https://goedge.cn .

package updates

import (
	"fmt"
	"github.com/TeaOSLab/EdgeAdmin/internal/utils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"os"
	"os/exec"
	"time"
)

var upgradeProgress float32
var isUpgrading = false

type UpgradeAction struct {
	actionutils.ParentAction
}

func (this *UpgradeAction) RunGet(params struct {
}) {
	this.Data["isUpgrading"] = isUpgrading
	this.Data["upgradeProgress"] = fmt.Sprintf("%.2f", upgradeProgress*100)
	this.Success()
}

func (this *UpgradeAction) RunPost(params struct {
	Url string
}) {
	if isUpgrading {
		this.Success()
		return
	}

	isUpgrading = true
	upgradeProgress = 0

	defer func() {
		isUpgrading = false
	}()

	var manager = utils.NewUpgradeManager("admin", params.Url)
	var ticker = time.NewTicker(1 * time.Second)
	go func() {
		for range ticker.C {
			if manager.IsDownloading() {
				var progress = manager.Progress()
				if progress >= 0 {
					upgradeProgress = progress
				}
			} else {
				return
			}
		}
	}()
	err := manager.Start()
	if err != nil {
		this.Fail("下载失败：" + err.Error())
		return
	}

	// restart
	exe, _ := os.Executable()
	if len(exe) > 0 {
		go func() {
			var cmd = exec.Command(exe, "restart")
			_ = cmd.Run()
		}()
	}

	this.Success()
}
