// Copyright 2021 Liuxiangchao iwind.liu@gmail.com. All rights reserved.

package logs

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/helpers"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
)

type FixAllAction struct {
	actionutils.ParentAction
}

func (this *FixAllAction) RunPost(params struct {
}) {
	defer this.CreateLogInfo("设置所有日志为已修复")

	_, err := this.RPC().NodeLogRPC().FixAllNodeLogs(this.AdminContext(), &pb.FixAllNodeLogsRequest{})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	// 通知左侧数字Badge更新
	helpers.NotifyNodeLogsCountChange()

	this.Success()
}
