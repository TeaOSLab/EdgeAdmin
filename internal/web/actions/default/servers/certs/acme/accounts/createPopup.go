// Copyright 2021 GoEdge CDN goedge.cdn@gmail.com. All rights reserved.

package accounts

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
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
	ProviderCode string
}) {
	this.Data["providerCode"] = params.ProviderCode

	// 服务商列表
	providersResp, err := this.RPC().ACMEProviderRPC().FindAllACMEProviders(this.AdminContext(), &pb.FindAllACMEProvidersRequest{})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	var providerMaps = []maps.Map{}
	for _, provider := range providersResp.AcmeProviders {
		providerMaps = append(providerMaps, maps.Map{
			"name":           provider.Name,
			"code":           provider.Code,
			"description":    provider.Description,
			"requireEAB":     provider.RequireEAB,
			"eabDescription": provider.EabDescription,
		})
	}

	this.Data["providers"] = providerMaps

	this.Show()
}

func (this *CreatePopupAction) RunPost(params struct {
	Name         string
	ProviderCode string
	EabKid       string
	EabKey       string

	Must *actions.Must
	CSRF *actionutils.CSRF
}) {
	var accountId int64
	defer func() {
		this.CreateLogInfo(codes.ACMEProviderAccount_LogCreateACMEProviderAccount, accountId)
	}()

	params.Must.
		Field("name", params.Name).
		Require("请输入账号名称").
		Field("providerCode", params.ProviderCode).
		Require("请选择服务商")

	providerResp, err := this.RPC().ACMEProviderRPC().FindACMEProviderWithCode(this.AdminContext(), &pb.FindACMEProviderWithCodeRequest{AcmeProviderCode: params.ProviderCode})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	var provider = providerResp.AcmeProvider
	if provider == nil {
		this.Fail("请选择服务商")
	}

	if provider.RequireEAB {
		params.Must.
			Field("eabKid", params.EabKid).
			Require("请输入EAB Kid").
			Field("eabKey", params.EabKey).
			Require("请输入EAB HMAC Key")
	}

	createResp, err := this.RPC().ACMEProviderAccountRPC().CreateACMEProviderAccount(this.AdminContext(), &pb.CreateACMEProviderAccountRequest{
		Name:         params.Name,
		ProviderCode: params.ProviderCode,
		EabKid:       params.EabKid,
		EabKey:       params.EabKey,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	accountId = createResp.AcmeProviderAccountId

	this.Success()
}
