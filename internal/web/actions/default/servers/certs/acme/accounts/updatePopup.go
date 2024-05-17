// Copyright 2021 GoEdge CDN goedge.cdn@gmail.com. All rights reserved.

package accounts

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
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
	AccountId int64
}) {
	// 账号信息
	accountResp, err := this.RPC().ACMEProviderAccountRPC().FindEnabledACMEProviderAccount(this.AdminContext(), &pb.FindEnabledACMEProviderAccountRequest{AcmeProviderAccountId: params.AccountId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	var account = accountResp.AcmeProviderAccount
	if account == nil {
		this.NotFound("ACMEProviderAccount", params.AccountId)
		return
	}

	var providerMap maps.Map
	if account.AcmeProvider != nil {
		providerMap = maps.Map{
			"name":           account.AcmeProvider.Name,
			"code":           account.AcmeProvider.Code,
			"description":    account.AcmeProvider.Description,
			"eabDescription": account.AcmeProvider.EabDescription,
			"requireEAB":     account.AcmeProvider.RequireEAB,
		}
	}

	this.Data["account"] = maps.Map{
		"id":           account.Id,
		"name":         account.Name,
		"isOn":         account.IsOn,
		"providerCode": account.ProviderCode,
		"eabKid":       account.EabKid,
		"eabKey":       account.EabKey,
		"provider":     providerMap,
	}

	this.Show()
}

func (this *UpdatePopupAction) RunPost(params struct {
	AccountId    int64
	Name         string
	ProviderCode string
	EabKid       string
	EabKey       string

	Must *actions.Must
	CSRF *actionutils.CSRF
}) {
	defer this.CreateLogInfo(codes.ACMEProviderAccount_LogUpdateACMEProviderAccount, params.AccountId)

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

	_, err = this.RPC().ACMEProviderAccountRPC().UpdateACMEProviderAccount(this.AdminContext(), &pb.UpdateACMEProviderAccountRequest{
		AcmeProviderAccountId: params.AccountId,
		Name:                  params.Name,
		EabKid:                params.EabKid,
		EabKey:                params.EabKey,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
