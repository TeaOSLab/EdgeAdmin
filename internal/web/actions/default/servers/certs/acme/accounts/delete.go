// Copyright 2021 GoEdge CDN goedge.cdn@gmail.com. All rights reserved.

package accounts

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
)

type DeleteAction struct {
	actionutils.ParentAction
}

func (this *DeleteAction) RunPost(params struct {
	AccountId int64
}) {
	defer this.CreateLogInfo(codes.ACMEProviderAccount_LogDeleteACMEProviderAccount, params.AccountId)

	_, err := this.RPC().ACMEProviderAccountRPC().DeleteACMEProviderAccount(this.AdminContext(), &pb.DeleteACMEProviderAccountRequest{AcmeProviderAccountId: params.AccountId})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
