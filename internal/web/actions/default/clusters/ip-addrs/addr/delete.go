// Copyright 2021 Liuxiangchao iwind.liu@gmail.com. All rights reserved.

package addr

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
)

type DeleteAction struct {
	actionutils.ParentAction
}

func (this *DeleteAction) RunPost(params struct {
	AddrId int64
}) {
	defer this.CreateLogInfo("删除IP地址 %d", params.AddrId)

	_, err := this.RPC().NodeIPAddressRPC().DisableNodeIPAddress(this.AdminContext(), &pb.DisableNodeIPAddressRequest{NodeIPAddressId: params.AddrId})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
