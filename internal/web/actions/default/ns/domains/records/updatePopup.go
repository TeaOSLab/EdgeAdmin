// Copyright 2021 Liuxiangchao iwind.liu@gmail.com. All rights reserved.

package records

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/dnsconfigs"
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
	RecordId int64
}) {
	recordResp, err := this.RPC().NSRecordRPC().FindEnabledNSRecord(this.AdminContext(), &pb.FindEnabledNSRecordRequest{NsRecordId: params.RecordId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	record := recordResp.NsRecord
	if record == nil {
		this.NotFound("nsRecord", params.RecordId)
		return
	}

	this.Data["record"] = maps.Map{
		"id":          record.Id,
		"name":        record.Name,
		"type":        record.Type,
		"value":       record.Value,
		"ttl":         record.Ttl,
		"weight":      record.Weight,
		"description": record.Description,
	}

	// 域名信息
	domainResp, err := this.RPC().NSDomainRPC().FindEnabledNSDomain(this.AdminContext(), &pb.FindEnabledNSDomainRequest{NsDomainId: record.NsDomain.Id})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	domain := domainResp.NsDomain
	if domain == nil {
		this.NotFound("nsDomain", record.NsDomain.Id)
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

func (this *UpdatePopupAction) RunPost(params struct {
	RecordId    int64
	Name        string
	Type        string
	Value       string
	Ttl         int32
	Description string

	Must *actions.Must
	CSRF *actionutils.CSRF
}) {
	this.CreateLogInfo("修改域名记录 %d", params.RecordId)

	_, err := this.RPC().NSRecordRPC().UpdateNSRecord(this.AdminContext(), &pb.UpdateNSRecordRequest{
		NsRecordId:  params.RecordId,
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

	this.Success()
}
