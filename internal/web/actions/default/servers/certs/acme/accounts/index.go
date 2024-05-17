// Copyright 2021 GoEdge CDN goedge.cdn@gmail.com. All rights reserved.

package accounts

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/maps"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.Nav("", "", "account")
}

func (this *IndexAction) RunGet(params struct{}) {
	countResp, err := this.RPC().ACMEProviderAccountRPC().CountAllEnabledACMEProviderAccounts(this.AdminContext(), &pb.CountAllEnabledACMEProviderAccountsRequest{})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	var count = countResp.Count
	var page = this.NewPage(count)
	this.Data["page"] = page.AsHTML()

	accountsResp, err := this.RPC().ACMEProviderAccountRPC().ListEnabledACMEProviderAccounts(this.AdminContext(), &pb.ListEnabledACMEProviderAccountsRequest{
		Offset: page.Offset,
		Size:   page.Size,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	var accountMaps = []maps.Map{}
	for _, account := range accountsResp.AcmeProviderAccounts {
		var providerMap maps.Map
		if account.AcmeProvider != nil {
			providerMap = maps.Map{
				"name":       account.AcmeProvider.Name,
				"code":       account.AcmeProvider.Code,
				"requireEAB": account.AcmeProvider.RequireEAB,
			}
		}

		accountMaps = append(accountMaps, maps.Map{
			"id":       account.Id,
			"isOn":     account.IsOn,
			"name":     account.Name,
			"eabKid":   account.EabKid,
			"eabKey":   account.EabKey,
			"provider": providerMap,
		})
	}
	this.Data["accounts"] = accountMaps

	this.Show()
}
