// Copyright 2021 GoEdge CDN goedge.cdn@gmail.com. All rights reserved.

package users

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/maps"
)

type AccountsWithCodeAction struct {
	actionutils.ParentAction
}

func (this *AccountsWithCodeAction) RunPost(params struct {
	Code string
}) {
	accountsResp, err := this.RPC().ACMEProviderAccountRPC().FindAllACMEProviderAccountsWithProviderCode(this.AdminContext(), &pb.FindAllACMEProviderAccountsWithProviderCodeRequest{AcmeProviderCode: params.Code})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	var accountMaps = []maps.Map{}
	for _, account := range accountsResp.AcmeProviderAccounts {
		accountMaps = append(accountMaps, maps.Map{
			"id":   account.Id,
			"name": account.Name,
		})
	}
	this.Data["accounts"] = accountMaps

	this.Success()
}
