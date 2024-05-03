// Copyright 2023 GoEdge CDN goedge.cdn@gmail.com. All rights reserved. Official site: https://goedge.cn .

package updates

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/TeaOSLab/EdgeAdmin/internal/utils"
	executils "github.com/TeaOSLab/EdgeAdmin/internal/utils/exec"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/nodeconfigs"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/maps"
	"os"
	"os/exec"
	"time"
)

var upgradeProgress float32
var isUpgrading = false
var isUpgradingDB = false

type UpgradeAction struct {
	actionutils.ParentAction
}

func (this *UpgradeAction) RunGet(params struct {
}) {
	this.Data["isUpgrading"] = isUpgrading
	this.Data["isUpgradingDB"] = isUpgradingDB
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

	// try to exec local 'edge-api upgrade'
	exePath, ok := this.checkLocalAPINode()
	if ok && len(exePath) > 0 {
		isUpgradingDB = true
		var before = time.Now()
		var cmd = executils.NewCmd(exePath, "upgrade")
		_ = cmd.Run()
		var costSeconds = time.Since(before).Seconds()

		// sleep to show upgrading status
		if costSeconds < 3 {
			time.Sleep(3 * time.Second)
		}
		isUpgradingDB = false
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

func (this *UpgradeAction) checkLocalAPINode() (exePath string, ok bool) {
	resp, err := this.RPC().APINodeRPC().FindCurrentAPINode(this.AdminContext(), &pb.FindCurrentAPINodeRequest{})
	if err != nil {
		return
	}
	if resp.ApiNode == nil {
		return
	}
	var instanceCode = resp.ApiNode.InstanceCode
	if len(instanceCode) == 0 {
		return
	}
	var statusJSON = resp.ApiNode.StatusJSON
	if len(statusJSON) == 0 {
		return
	}

	var status = &nodeconfigs.NodeStatus{}
	err = json.Unmarshal(statusJSON, status)
	if err != nil {
		return
	}

	exePath = status.ExePath
	if len(exePath) == 0 {
		return
	}

	stat, err := os.Stat(exePath)
	if err != nil {
		return
	}
	if stat.IsDir() {
		return
	}

	// 实例信息
	{
		var outputBuffer = &bytes.Buffer{}
		var cmd = exec.Command(exePath, "instance")
		cmd.Stdout = outputBuffer
		err = cmd.Run()
		if err != nil {
			return
		}

		var outputBytes = outputBuffer.Bytes()
		if len(outputBytes) == 0 {
			return
		}

		var instanceMap = maps.Map{}
		err = json.Unmarshal(bytes.TrimSpace(outputBytes), &instanceMap)
		if err != nil {
			return
		}

		if instanceMap.GetString("code") != instanceCode {
			return
		}
	}

	ok = true
	return
}
