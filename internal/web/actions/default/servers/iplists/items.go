// Copyright 2021 Liuxiangchao iwind.liu@gmail.com. All rights reserved.

package iplists

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs/firewallconfigs"
	"github.com/iwind/TeaGo/maps"
	timeutil "github.com/iwind/TeaGo/utils/time"
)

type ItemsAction struct {
	actionutils.ParentAction
}

func (this *ItemsAction) Init() {
	this.Nav("", "", "item")
}

func (this *ItemsAction) RunGet(params struct {
	ListId int64
}) {
	err := InitIPList(this.Parent(), params.ListId)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	// 数量
	var listId = params.ListId
	countResp, err := this.RPC().IPItemRPC().CountIPItemsWithListId(this.AdminContext(), &pb.CountIPItemsWithListIdRequest{IpListId: listId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	count := countResp.Count
	page := this.NewPage(count)
	this.Data["page"] = page.AsHTML()

	// 列表
	itemsResp, err := this.RPC().IPItemRPC().ListIPItemsWithListId(this.AdminContext(), &pb.ListIPItemsWithListIdRequest{
		IpListId: listId,
		Offset:   page.Offset,
		Size:     page.Size,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	itemMaps := []maps.Map{}
	for _, item := range itemsResp.IpItems {
		expiredTime := ""
		if item.ExpiredAt > 0 {
			expiredTime = timeutil.FormatTime("Y-m-d H:i:s", item.ExpiredAt)
		}

		itemMaps = append(itemMaps, maps.Map{
			"id":             item.Id,
			"ipFrom":         item.IpFrom,
			"ipTo":           item.IpTo,
			"expiredTime":    expiredTime,
			"reason":         item.Reason,
			"type":           item.Type,
			"eventLevelName": firewallconfigs.FindFirewallEventLevelName(item.EventLevel),
		})
	}
	this.Data["items"] = itemMaps

	this.Show()
}
