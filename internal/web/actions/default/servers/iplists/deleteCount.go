// Copyright 2023 GoEdge CDN goedge.cdn@gmail.com. All rights reserved.

package iplists

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/helpers"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/types"
	"strings"
)

type DeleteCountAction struct {
	actionutils.ParentAction
}

func (this *DeleteCountAction) RunPost(params struct {
	Ip         string
	Keyword    string
	GlobalOnly bool
	Unread     bool
	EventLevel string
	ListType   string

	Count int64
}) {

	var count = params.Count
	if count <= 0 || count >= 100_000 {
		this.Fail("'count' 参数错误")
		return
	}

	itemIdsResp, err := this.RPC().IPItemRPC().ListAllIPItemIds(this.AdminContext(), &pb.ListAllIPItemIdsRequest{
		Keyword:    params.Keyword,
		GlobalOnly: params.GlobalOnly,
		Unread:     params.Unread,
		EventLevel: params.EventLevel,
		ListType:   params.ListType,
		Ip:         params.Ip,
		Offset:     0,
		Size:       count,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	var itemIds = itemIdsResp.IpItemIds

	if len(itemIds) == 0 {
		this.Success()
	}

	// 记录日志
	defer func() {
		var itemIdStrings = []string{}
		for _, itemId := range itemIds {
			itemIdStrings = append(itemIdStrings, types.String(itemId))
		}

		var itemIdsDescription = ""
		if len(itemIdStrings) > 10 {
			itemIdsDescription = strings.Join(itemIdStrings[:10], ", ") + " ..."
		} else {
			itemIdsDescription = strings.Join(itemIdStrings, ", ")
		}
		this.CreateLogInfo(codes.IPList_LogDeleteIPBatch, itemIdsDescription)
	}()

	_, err = this.RPC().IPItemRPC().DeleteIPItems(this.AdminContext(), &pb.DeleteIPItemsRequest{IpItemIds: itemIds})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	// 通知左侧菜单Badge更新
	helpers.NotifyIPItemsCountChanges()

	this.Success()
}
