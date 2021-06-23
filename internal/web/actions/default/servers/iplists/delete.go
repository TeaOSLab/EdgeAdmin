// Copyright 2021 Liuxiangchao iwind.liu@gmail.com. All rights reserved.

package iplists

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
)

type DeleteAction struct {
	actionutils.ParentAction
}

func (this *DeleteAction) RunPost(params struct {
	ListId int64
}) {
	defer this.CreateLogInfo("删除IP名单 %d", params.ListId)

	// 删除
	_, err := this.RPC().IPListRPC().DeleteIPList(this.AdminContext(), &pb.DeleteIPListRequest{IpListId: params.ListId})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
