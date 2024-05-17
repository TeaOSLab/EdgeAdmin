// Copyright 2022 GoEdge CDN goedge.cdn@gmail.com. All rights reserved. Official site: https://goedge.cn .

package dashboard

import (
	"bytes"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"os/exec"
	"regexp"
	"time"
)

type RestartLocalAPINodeAction struct {
	actionutils.ParentAction
}

func (this *RestartLocalAPINodeAction) RunPost(params struct {
	ExePath string
}) {
	// 检查当前用户是超级用户
	adminResp, err := this.RPC().AdminRPC().FindEnabledAdmin(this.AdminContext(), &pb.FindEnabledAdminRequest{AdminId: this.AdminId()})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	if adminResp.Admin == nil || !adminResp.Admin.IsSuper {
		this.Fail("请切换到超级用户进行此操作")
	}

	var exePath = params.ExePath
	if len(exePath) == 0 {
		this.Fail("找不到要重启的API节点文件")
	}

	{
		var stdoutBuffer = &bytes.Buffer{}
		var cmd = exec.Command(exePath, "restart")
		cmd.Stdout = stdoutBuffer
		err = cmd.Run()
		if err != nil {
			this.Fail("运行失败：输出：" + stdoutBuffer.String())
		}
	}

	// 停止1秒等待命令运行完毕
	time.Sleep(1 * time.Second)

	// 检查是否已启动
	var countTries = 120
	for {
		countTries--
		if countTries < 0 {
			this.Fail("启动超时，请尝试手动启动")
			break
		}

		var stdoutBuffer = &bytes.Buffer{}
		var cmd = exec.Command(exePath, "status")
		cmd.Stdout = stdoutBuffer
		err = cmd.Run()
		if err != nil {
			time.Sleep(1 * time.Second)
			continue
		}

		if regexp.MustCompile(`pid:\s*\d+`).
			MatchString(stdoutBuffer.String()) {
			break
		}

		time.Sleep(1 * time.Second)
	}

	this.Success()
}
