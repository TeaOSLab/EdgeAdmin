// Copyright 2021 Liuxiangchao iwind.liu@gmail.com. All rights reserved.

package clusters

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/actions"
)

type CreatePopupAction struct {
	actionutils.ParentAction
}

func (this *CreatePopupAction) Init() {
}

func (this *CreatePopupAction) RunGet(params struct {
	ClusterId int64
	DomainId  int64
	UserId    int64
}) {
	this.Data["clusterId"] = params.ClusterId
	this.Data["domainId"] = params.DomainId
	this.Data["userId"] = params.UserId

	this.Show()
}

func (this *CreatePopupAction) RunPost(params struct {
	ClusterId int64
	DomainId  int64
	UserId    int64

	Name       string
	RangesJSON []byte

	Must *actions.Must
	CSRF *actionutils.CSRF
}) {
	var routeId = int64(0)
	defer func() {
		this.CreateLogInfo("创建域名服务线路 %d", routeId)
	}()

	params.Must.Field("name", params.Name).
		Require("请输入线路名称")

	createResp, err := this.RPC().NSRouteRPC().CreateNSRoute(this.AdminContext(), &pb.CreateNSRouteRequest{
		NsClusterId: params.ClusterId,
		NsDomainId:  params.DomainId,
		UserId:      params.UserId,
		Name:        params.Name,
		RangesJSON:  params.RangesJSON,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	routeId = createResp.NsRouteId

	this.Success()
}
