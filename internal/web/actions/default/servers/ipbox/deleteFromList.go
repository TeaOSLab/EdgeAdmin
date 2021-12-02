// Copyright 2021 Liuxiangchao iwind.liu@gmail.com. All rights reserved.

package ipbox

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
)

type DeleteFromListAction struct {
	actionutils.ParentAction
}

func (this *DeleteFromListAction) RunPost(params struct {
	ListId int64
	ItemId int64
}) {
	defer this.CreateLogInfo("从IP名单 %d 中删除IP %d", params.ListId, params.ItemId)

	_, err := this.RPC().IPItemRPC().DeleteIPItem(this.AdminContext(), &pb.DeleteIPItemRequest{IpItemId: params.ItemId})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
