// Copyright 2021 Liuxiangchao iwind.liu@gmail.com. All rights reserved.

package records

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/dnsconfigs"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type CreatePopupAction struct {
	actionutils.ParentAction
}

func (this *CreatePopupAction) Init() {
	this.Nav("", "", "")
}

func (this *CreatePopupAction) RunGet(params struct {
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
	this.Data["domain"] = maps.Map{
		"id":   domain.Id,
		"name": domain.Name,
	}

	// 类型
	this.Data["types"] = dnsconfigs.FindAllRecordTypeDefinitions()

	// TTL
	this.Data["ttlValues"] = dnsconfigs.FindAllRecordTTL()

	this.Show()
}

func (this *CreatePopupAction) RunPost(params struct {
	DomainId    int64
	Name        string
	Type        string
	Value       string
	Ttl         int32
	Description string

	Must *actions.Must
	CSRF *actionutils.CSRF
}) {
	var recordId int64
	defer func() {
		this.CreateLogInfo("创建域名记录 %d", recordId)
	}()

	createResp, err := this.RPC().NSRecordRPC().CreateNSRecord(this.AdminContext(), &pb.CreateNSRecordRequest{
		NsDomainId:  params.DomainId,
		Description: params.Description,
		Name:        params.Name,
		Type:        params.Type,
		Value:       params.Value,
		Ttl:         params.Ttl,
		NsRouteIds:  nil, // TODO 等待实现
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	recordId = createResp.NsRecordId

	this.Success()
}
