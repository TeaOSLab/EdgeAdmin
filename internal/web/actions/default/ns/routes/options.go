// Copyright 2021 Liuxiangchao iwind.liu@gmail.com. All rights reserved.

package clusters

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/maps"
)

type OptionsAction struct {
	actionutils.ParentAction
}

func (this *OptionsAction) RunPost(params struct {
	ClusterId int64
	DomainId  int64
	UserId    int64
}) {
	routesResp, err := this.RPC().NSRouteRPC().FindAllEnabledNSRoutes(this.AdminContext(), &pb.FindAllEnabledNSRoutesRequest{
		NsClusterId: params.ClusterId,
		NsDomainId:  params.DomainId,
		UserId:      params.UserId,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	routeMaps := []maps.Map{}
	for _, route := range routesResp.NsRoutes {
		routeMaps = append(routeMaps, maps.Map{
			"id":   route.Id,
			"name": route.Name,
		})
	}
	this.Data["routes"] = routeMaps

	this.Success()
}
