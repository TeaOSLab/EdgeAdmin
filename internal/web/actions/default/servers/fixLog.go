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
	LogId int64
}) {
	_, err := this.RPC().NodeLogRPC().FixNodeLog(this.AdminContext(), &pb.FixNodeLogRequest{NodeLogId: params.LogId})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
