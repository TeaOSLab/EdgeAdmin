package iplists

import (	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/helpers"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
)

type DeleteIPAction struct {
	actionutils.ParentAction
}

func (this *DeleteIPAction) RunPost(params struct {
	ItemId int64
}) {
	// 日志
	defer this.CreateLogInfo(codes.IPItem_LogDeleteIPItem, params.ItemId)

	_, err := this.RPC().IPItemRPC().DeleteIPItem(this.AdminContext(), &pb.DeleteIPItemRequest{IpItemId: params.ItemId})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	// 通知左侧菜单Badge更新
	helpers.NotifyIPItemsCountChanges()

	this.Success()
}
