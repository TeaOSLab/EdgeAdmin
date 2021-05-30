// Copyright 2021 Liuxiangchao iwind.liu@gmail.com. All rights reserved.

package clusters

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
)

type DeleteAction struct {
	actionutils.ParentAction
}

func (this *DeleteAction) RunPost(params struct {
	RouteId int64
}) {
	defer this.CreateLogInfo("删除域名服务线路 %d", params.RouteId)

	_, err := this.RPC().NSRouteRPC().DeleteNSRoute(this.AdminContext(), &pb.DeleteNSRouteRequest{NsRouteId: params.RouteId})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
