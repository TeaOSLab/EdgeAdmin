// Copyright 2021 Liuxiangchao iwind.liu@gmail.com. All rights reserved.

package records

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/dns/domains/domainutils"
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
	RouteCodes  []string

	Must *actions.Must
	CSRF *actionutils.CSRF
}) {
	var recordId int64
	defer func() {
		this.CreateLogInfo("创建域名记录 %d", recordId)
	}()

	// 校验记录名
	if !domainutils.ValidateRecordName(params.Name) {
		this.FailField("name", "请输入正确的记录名")
	}

	// 校验记录值
	message, ok := domainutils.ValidateRecordValue(params.Type, params.Value)
	if !ok {
		this.FailField("value", "记录值错误："+message)
	}

	createResp, err := this.RPC().NSRecordRPC().CreateNSRecord(this.AdminContext(), &pb.CreateNSRecordRequest{
		NsDomainId:   params.DomainId,
		Description:  params.Description,
		Name:         params.Name,
		Type:         params.Type,
		Value:        params.Value,
		Ttl:          params.Ttl,
		NsRouteCodes: params.RouteCodes,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	recordId = createResp.NsRecordId

	this.Success()
}
