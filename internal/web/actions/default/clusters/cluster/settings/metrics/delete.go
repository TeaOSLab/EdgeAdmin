// Copyright 2021 Liuxiangchao iwind.liu@gmail.com. All rights reserved.

package metrics

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
)

type DeleteAction struct {
	actionutils.ParentAction
}

func (this *DeleteAction) RunPost(params struct {
	ClusterId int64
	ItemId    int64
}) {
	defer this.CreateLogInfo("从集群 %d 中移除指标 %d", params.ClusterId, params.ItemId)

	_, err := this.RPC().NodeClusterMetricItemRPC().DisableNodeClusterMetricItem(this.AdminContext(), &pb.DisableNodeClusterMetricItemRequest{
		NodeClusterId: params.ClusterId,
		MetricItemId:  params.ItemId,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
