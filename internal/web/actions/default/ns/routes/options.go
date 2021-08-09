// Copyright 2021 Liuxiangchao iwind.liu@gmail.com. All rights reserved.

package clusters

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/dnsconfigs"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/types"
)

type OptionsAction struct {
	actionutils.ParentAction
}

func (this *OptionsAction) RunPost(params struct {
	ClusterId int64
	DomainId  int64
	UserId    int64
}) {
	var routeMaps = []maps.Map{}

	// 默认线路
	for _, route := range dnsconfigs.AllDefaultRoutes {
		routeMaps = append(routeMaps, maps.Map{
			"name": route.Name,
			"code": route.Code,
			"type": "default",
		})
	}

	// 自定义
	routesResp, err := this.RPC().NSRouteRPC().FindAllEnabledNSRoutes(this.AdminContext(), &pb.FindAllEnabledNSRoutesRequest{
		NsClusterId: params.ClusterId,
		NsDomainId:  params.DomainId,
		UserId:      params.UserId,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	for _, route := range routesResp.NsRoutes {
		if len(route.Code) == 0 {
			route.Code = "id:" + types.String(route.Id)
		}

		routeMaps = append(routeMaps, maps.Map{
			"name": route.Name,
			"code": route.Code,
			"type": "user",
		})
	}

	// 运营商
	for _, route := range dnsconfigs.AllDefaultISPRoutes {
		routeMaps = append(routeMaps, maps.Map{
			"name": route.Name,
			"code": route.Code,
			"type": "isp",
		})
	}

	// 中国
	for _, route := range dnsconfigs.AllDefaultChinaProvinceRoutes {
		routeMaps = append(routeMaps, maps.Map{
			"name": route.Name,
			"code": route.Code,
			"type": "china",
		})
	}

	// 全球
	for _, route := range dnsconfigs.AllDefaultWorldRegionRoutes {
		routeMaps = append(routeMaps, maps.Map{
			"name": route.Name,
			"code": route.Code,
			"type": "world",
		})
	}

	this.Data["routes"] = routeMaps

	this.Success()
}
