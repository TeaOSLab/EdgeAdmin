// Copyright 2021 Liuxiangchao iwind.liu@gmail.com. All rights reserved.

package domains

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/ns/domains/domainutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/maps"
)

type DomainAction struct {
	actionutils.ParentAction
}

func (this *DomainAction) Init() {
	this.Nav("", "", "index")
}

func (this *DomainAction) RunGet(params struct {
	DomainId int64
}) {
	err := domainutils.InitDomain(this.Parent(), params.DomainId)
	if err != nil {
		this.ErrorPage(err)
		return
	}
	var countRecords = this.Data.GetMap("domain").GetInt64("countRecords")
	var countKeys = this.Data.GetMap("domain").GetInt64("countKeys")

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

	var clusterMap maps.Map
	if domain.NsCluster != nil {
		clusterMap = maps.Map{
			"id":   domain.NsCluster.Id,
			"name": domain.NsCluster.Name,
		}
	}

	// 用户信息
	var userMap maps.Map
	if domain.User != nil {
		userMap = maps.Map{
			"id":       domain.User.Id,
			"username": domain.User.Username,
			"fullname": domain.User.Fullname,
		}
	}

	this.Data["domain"] = maps.Map{
		"id":           domain.Id,
		"name":         domain.Name,
		"isOn":         domain.IsOn,
		"cluster":      clusterMap,
		"user":         userMap,
		"countRecords": countRecords,
		"countKeys":    countKeys,
	}

	this.Show()
}
