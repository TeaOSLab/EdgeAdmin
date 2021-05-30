// Copyright 2021 Liuxiangchao iwind.liu@gmail.com. All rights reserved.

package clusters

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/maps"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.Nav("", "", "index")
}

func (this *IndexAction) RunGet(params struct{}) {
	routesResp, err := this.RPC().NSRouteRPC().FindAllEnabledNSRoutes(this.AdminContext(), &pb.FindAllEnabledNSRoutesRequest{
		NsClusterId: 0,
		NsDomainId:  0,
		UserId:      0,
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
			"isOn": route.IsOn,
		})
	}
	this.Data["routes"] = routeMaps

	this.Show()
}
