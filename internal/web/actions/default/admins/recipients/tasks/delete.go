// Copyright 2021 Liuxiangchao iwind.liu@gmail.com. All rights reserved.

package tasks

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
)

type DeleteAction struct {
	actionutils.ParentAction
}

func (this *DeleteAction) RunPost(params struct {
	TaskId int64
}) {
	defer this.CreateLogInfo("删除消息发送任务 %d", params.TaskId)

	_, err := this.RPC().MessageTaskRPC().DeleteMessageTask(this.AdminContext(), &pb.DeleteMessageTaskRequest{MessageTaskId: params.TaskId})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
