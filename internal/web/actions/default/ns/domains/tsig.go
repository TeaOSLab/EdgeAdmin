// Copyright 2021 Liuxiangchao iwind.liu@gmail.com. All rights reserved.

package domains

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/ns/domains/domainutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/dnsconfigs"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/logs"
)

type TsigAction struct {
	actionutils.ParentAction
}

func (this *TsigAction) Init() {
	this.Nav("", "", "tsig")
}

func (this *TsigAction) RunGet(params struct {
	DomainId int64
}) {
	// 初始化域名信息
	err := domainutils.InitDomain(this.Parent(), params.DomainId)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	// TSIG信息
	tsigResp, err := this.RPC().NSDomainRPC().FindEnabledNSDomainTSIG(this.AdminContext(), &pb.FindEnabledNSDomainTSIGRequest{NsDomainId: params.DomainId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	var tsigJSON = tsigResp.TsigJSON

	var tsigConfig = &dnsconfigs.TSIGConfig{}
	if len(tsigJSON) > 0 {
		err = json.Unmarshal(tsigJSON, tsigConfig)
		if err != nil {
			// 只是提示错误，仍然允许用户修改
			logs.Error(err)
		}
	}
	this.Data["tsig"] = tsigConfig

	this.Show()
}

func (this *TsigAction) RunPost(params struct {
	DomainId int64
	IsOn     bool

	Must *actions.Must
	CSRF *actionutils.CSRF
}) {
	defer this.CreateLogInfo("修改DNS域名 %d 的TSIG配置", params.DomainId)

	var tsigConfig = &dnsconfigs.TSIGConfig{
		IsOn: params.IsOn,
	}
	tsigJSON, err := json.Marshal(tsigConfig)
	if err != nil {
		this.ErrorPage(err)
		return
	}
	_, err = this.RPC().NSDomainRPC().UpdateNSDomainTSIG(this.AdminContext(), &pb.UpdateNSDomainTSIGRequest{
		NsDomainId: params.DomainId,
		TsigJSON:   tsigJSON,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	this.Success()
}
