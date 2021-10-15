// Copyright 2021 Liuxiangchao iwind.liu@gmail.com. All rights reserved.

package logs

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
)

type ReadLogsAction struct {
	actionutils.ParentAction
}

func (this *ReadLogsAction) RunPost(params struct {
	LogIds []int64
}) {
	_, err := this.RPC().NodeLogRPC().UpdateNodeLogsRead(this.AdminContext(), &pb.UpdateNodeLogsReadRequest{
		NodeLogIds: params.LogIds,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
