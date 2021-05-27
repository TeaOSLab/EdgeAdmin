// Copyright 2021 Liuxiangchao iwind.liu@gmail.com. All rights reserved.

package domains

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
)

type DeleteAction struct {
	actionutils.ParentAction
}

func (this *DeleteAction) RunPost(params struct {
	DomainId int64
}) {
	defer this.CreateLogInfo("删除域名 %d", params.DomainId)

	_, err := this.RPC().NSDomainRPC().DeleteNSDomain(this.AdminContext(), &pb.DeleteNSDomainRequest{NsDomainId: params.DomainId})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
