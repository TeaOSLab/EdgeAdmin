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
	ListId  int64
	Keyword string
}) {
	this.Data["keyword"] = params.Keyword

	err := InitIPList(this.Parent(), params.ListId)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	// 数量
	var listId = params.ListId
	countResp, err := this.RPC().IPItemRPC().CountIPItemsWithListId(this.AdminContext(), &pb.CountIPItemsWithListIdRequest{
		IpListId: listId,
		Keyword:  params.Keyword,
	})
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
		Keyword:  params.Keyword,
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

		// policy
		var sourcePolicyMap = maps.Map{"id": 0}
		if item.SourceHTTPFirewallPolicy != nil {
			sourcePolicyMap = maps.Map{
				"id":       item.SourceHTTPFirewallPolicy.Id,
				"name":     item.SourceHTTPFirewallPolicy.Name,
				"serverId": item.SourceHTTPFirewallPolicy.ServerId,
			}
		}

		// group
		var sourceGroupMap = maps.Map{"id": 0}
		if item.SourceHTTPFirewallRuleGroup != nil {
			sourceGroupMap = maps.Map{
				"id":   item.SourceHTTPFirewallRuleGroup.Id,
				"name": item.SourceHTTPFirewallRuleGroup.Name,
			}
		}

		// set
		var sourceSetMap = maps.Map{"id": 0}
		if item.SourceHTTPFirewallRuleSet != nil {
			sourceSetMap = maps.Map{
				"id":   item.SourceHTTPFirewallRuleSet.Id,
				"name": item.SourceHTTPFirewallRuleSet.Name,
			}
		}

		// server
		var sourceServerMap = maps.Map{"id": 0}
		if item.SourceServer != nil {
			sourceServerMap = maps.Map{
				"id":   item.SourceServer.Id,
				"name": item.SourceServer.Name,
			}
		}

		itemMaps = append(itemMaps, maps.Map{
			"id":             item.Id,
			"ipFrom":         item.IpFrom,
			"ipTo":           item.IpTo,
			"createdTime":    timeutil.FormatTime("Y-m-d", item.CreatedAt),
			"expiredTime":    expiredTime,
			"reason":         item.Reason,
			"type":           item.Type,
			"eventLevelName": firewallconfigs.FindFirewallEventLevelName(item.EventLevel),
			"sourcePolicy":   sourcePolicyMap,
			"sourceGroup":    sourceGroupMap,
			"sourceSet":      sourceSetMap,
			"sourceServer":   sourceServerMap,
		})
	}
	this.Data["items"] = itemMaps

	this.Show()
}
