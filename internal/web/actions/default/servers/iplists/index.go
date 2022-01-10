// Copyright 2021 Liuxiangchao iwind.liu@gmail.com. All rights reserved.

package iplists

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs/firewallconfigs"
	"github.com/iwind/TeaGo/maps"
	timeutil "github.com/iwind/TeaGo/utils/time"
	"time"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.Nav("", "", "index")
}

func (this *IndexAction) RunGet(params struct {
	Ip         string
	GlobalOnly bool
	Unread     bool
}) {
	this.Data["type"] = ""
	this.Data["ip"] = params.Ip
	this.Data["globalOnly"] = params.GlobalOnly
	this.Data["unread"] = params.Unread

	countUnreadResp, err := this.RPC().IPItemRPC().CountAllEnabledIPItems(this.AdminContext(), &pb.CountAllEnabledIPItemsRequest{
		Ip:         params.Ip,
		GlobalOnly: params.GlobalOnly,
		Unread:     true,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	this.Data["countUnread"] = countUnreadResp.Count

	countResp, err := this.RPC().IPItemRPC().CountAllEnabledIPItems(this.AdminContext(), &pb.CountAllEnabledIPItemsRequest{
		Ip:         params.Ip,
		GlobalOnly: params.GlobalOnly,
		Unread:     params.Unread,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	var count = countResp.Count
	var page = this.NewPage(count)
	this.Data["page"] = page.AsHTML()

	itemsResp, err := this.RPC().IPItemRPC().ListAllEnabledIPItems(this.AdminContext(), &pb.ListAllEnabledIPItemsRequest{
		Ip:         params.Ip,
		GlobalOnly: params.GlobalOnly,
		Unread:     params.Unread,
		Offset:     page.Offset,
		Size:       page.Size,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	var itemMaps = []maps.Map{}
	for _, result := range itemsResp.Results {
		var item = result.IpItem
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

		// IP名单
		var listMap = maps.Map{"id": 0}
		if result.IpList != nil {
			listMap = maps.Map{
				"id":   result.IpList.Id,
				"name": result.IpList.Name,
				"type": result.IpList.Type,
			}
		}

		// policy
		var policyMap = maps.Map{"id": 0}
		if result.HttpFirewallPolicy != nil {
			policyMap = maps.Map{
				"id":   result.HttpFirewallPolicy.Id,
				"name": result.HttpFirewallPolicy.Name,
			}

			if result.Server != nil {
				policyMap["server"] = maps.Map{"id": result.Server.Id, "name": result.Server.Name}
			}
		}

		// node
		var sourceNodeMap = maps.Map{"id": 0}
		if item.SourceNode != nil && item.SourceNode.NodeCluster != nil {
			sourceNodeMap = maps.Map{
				"id":        item.SourceNode.Id,
				"name":      item.SourceNode.Name,
				"clusterId": item.SourceNode.NodeCluster.Id,
			}
		}

		itemMaps = append(itemMaps, maps.Map{
			"id":             item.Id,
			"ipFrom":         item.IpFrom,
			"ipTo":           item.IpTo,
			"createdTime":    timeutil.FormatTime("Y-m-d", item.CreatedAt),
			"isExpired":      item.ExpiredAt < time.Now().Unix(),
			"expiredTime":    expiredTime,
			"reason":         item.Reason,
			"type":           item.Type,
			"isRead":         item.IsRead,
			"lifeSeconds":    item.ExpiredAt - time.Now().Unix(),
			"eventLevelName": firewallconfigs.FindFirewallEventLevelName(item.EventLevel),
			"sourcePolicy":   sourcePolicyMap,
			"sourceGroup":    sourceGroupMap,
			"sourceSet":      sourceSetMap,
			"sourceServer":   sourceServerMap,
			"sourceNode":     sourceNodeMap,
			"list":           listMap,
			"policy":         policyMap,
		})
	}
	this.Data["items"] = itemMaps

	this.Show()
}
