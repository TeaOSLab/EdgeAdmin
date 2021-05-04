// Copyright 2021 Liuxiangchao iwind.liu@gmail.com. All rights reserved.

package thresholds

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
)

// DeleteAction 删除阈值
type DeleteAction struct {
	actionutils.ParentAction
}

func (this *DeleteAction) RunPost(params struct {
	ThresholdId int64
}) {
	defer this.CreateLogInfo("删除阈值 %d", params.ThresholdId)

	_, err := this.RPC().NodeThresholdRPC().DeleteNodeThreshold(this.AdminContext(), &pb.DeleteNodeThresholdRequest{NodeThresholdId: params.ThresholdId})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
