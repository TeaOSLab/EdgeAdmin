// Copyright 2021 Liuxiangchao iwind.liu@gmail.com. All rights reserved.

package clusters

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
)

type SortAction struct {
	actionutils.ParentAction
}

func (this *SortAction) RunPost(params struct {
	RouteIds []int64
}) {
	defer this.CreateLogInfo("对线路进行排序")

	_, err := this.RPC().NSRouteRPC().UpdateNSRouteOrders(this.AdminContext(), &pb.UpdateNSRouteOrdersRequest{NsRouteIds: params.RouteIds})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
