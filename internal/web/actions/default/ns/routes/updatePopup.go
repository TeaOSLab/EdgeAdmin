// Copyright 2021 Liuxiangchao iwind.liu@gmail.com. All rights reserved.

package clusters

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type UpdatePopupAction struct {
	actionutils.ParentAction
}

func (this *UpdatePopupAction) Init() {
	this.Nav("", "", "")
}

func (this *UpdatePopupAction) RunGet(params struct {
	RouteId int64
}) {
	routeResp, err := this.RPC().NSRouteRPC().FindEnabledNSRoute(this.AdminContext(), &pb.FindEnabledNSRouteRequest{NsRouteId: params.RouteId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	route := routeResp.NsRoute
	if route == nil {
		this.NotFound("nsRoute", params.RouteId)
		return
	}

	rangeMaps := []maps.Map{}
	if len(route.RangesJSON) > 0 {
		err = json.Unmarshal([]byte(route.RangesJSON), &rangeMaps)
		if err != nil {
			this.ErrorPage(err)
			return
		}
	}

	this.Data["route"] = maps.Map{
		"id":     route.Id,
		"name":   route.Name,
		"isOn":   route.IsOn,
		"ranges": rangeMaps,
	}

	this.Show()
}

func (this *UpdatePopupAction) RunPost(params struct {
	RouteId    int64
	Name       string
	RangesJSON []byte

	Must *actions.Must
	CSRF *actionutils.CSRF
}) {
	defer this.CreateLogInfo("修改域名线路 %d", params.RouteId)

	params.Must.Field("name", params.Name).
		Require("请输入线路名称")

	_, err := this.RPC().NSRouteRPC().UpdateNSRoute(this.AdminContext(), &pb.UpdateNSRouteRequest{
		NsRouteId:  params.RouteId,
		Name:       params.Name,
		RangesJSON: params.RangesJSON,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
