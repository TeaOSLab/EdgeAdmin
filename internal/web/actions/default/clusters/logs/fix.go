// Copyright 2021 Liuxiangchao iwind.liu@gmail.com. All rights reserved.

package logs

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/helpers"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/types"
	"strings"
)

type FixAction struct {
	actionutils.ParentAction
}

func (this *FixAction) RunPost(params struct {
	LogIds []int64
}) {
	var logIdStrings = []string{}
	for _, logId := range params.LogIds {
		logIdStrings = append(logIdStrings, types.String(logId))
	}

	defer this.CreateLogInfo("设置日志 %s 为已修复", strings.Join(logIdStrings, ", "))

	_, err := this.RPC().NodeLogRPC().FixNodeLogs(this.AdminContext(), &pb.FixNodeLogsRequest{NodeLogIds: params.LogIds})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	// 通知左侧数字Badge更新
	helpers.NotifyNodeLogsCountChange()

	this.Success()
}
