// Copyright 2021 Liuxiangchao iwind.liu@gmail.com. All rights reserved.

package domains

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/dns/domains/domainutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type UpdateAction struct {
	actionutils.ParentAction
}

func (this *UpdateAction) Init() {
	this.Nav("", "", "update")
}

func (this *UpdateAction) RunGet(params struct {
	DomainId int64
}) {
	// 域名信息
	domainResp, err := this.RPC().NSDomainRPC().FindEnabledNSDomain(this.AdminContext(), &pb.FindEnabledNSDomainRequest{NsDomainId: params.DomainId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	domain := domainResp.NsDomain
	if domain == nil {
		this.NotFound("nsDomain", params.DomainId)
		return
	}

	var clusterId = int64(0)
	if domain.NsCluster != nil {
		clusterId = domain.NsCluster.Id
	}

	// 用户信息
	var userId = int64(0)
	if domain.User != nil {
		userId = domain.User.Id
	}

	this.Data["domain"] = maps.Map{
		"id":        domain.Id,
		"name":      domain.Name,
		"isOn":      domain.IsOn,
		"clusterId": clusterId,
		"userId":    userId,
	}

	this.Show()
}

func (this *UpdateAction) RunPost(params struct {
	DomainId  int64
	Name      string
	ClusterId int64
	UserId    int64
	IsOn      bool

	Must *actions.Must
	CSRF *actionutils.CSRF
}) {
	this.CreateLogInfo("修改域名 %d", params.DomainId)

	params.Must.
		Field("name", params.Name).
		Require("请输入域名").
		Expect(func() (message string, success bool) {
			success = domainutils.ValidateDomainFormat(params.Name)
			if !success {
				message = "请输入正确的域名"
			}
			return
		}).
		Field("clusterId", params.ClusterId).
		Gt(0, "请选择所属集群")

	_, err := this.RPC().NSDomainRPC().UpdateNSDomain(this.AdminContext(), &pb.UpdateNSDomainRequest{
		NsDomainId:  params.DomainId,
		NsClusterId: params.ClusterId,
		UserId:      params.UserId,
		Name:        params.Name,
		IsOn:        params.IsOn,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
