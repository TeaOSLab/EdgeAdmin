// Copyright 2021 Liuxiangchao iwind.liu@gmail.com. All rights reserved.

package servers

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
)

type FixLogAction struct {
	actionutils.ParentAction
}

func (this *FixLogAction) RunPost(params struct {
	LogIds []int64
}) {
	_, err := this.RPC().NodeLogRPC().FixNodeLogs(this.AdminContext(), &pb.FixNodeLogsRequest{NodeLogIds: params.LogIds})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
