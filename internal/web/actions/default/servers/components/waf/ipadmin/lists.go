package ipadmin

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/dao"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs/firewallconfigs"
	"github.com/iwind/TeaGo/maps"
	timeutil "github.com/iwind/TeaGo/utils/time"
)

type ListsAction struct {
	actionutils.ParentAction
}

func (this *ListsAction) Init() {
	this.Nav("", "", "ipadmin")
}

func (this *ListsAction) RunGet(params struct {
	FirewallPolicyId int64
	Type             string
}) {
	this.Data["subMenuItem"] = params.Type
	this.Data["type"] = params.Type

	listId, err := dao.SharedHTTPFirewallPolicyDAO.FindEnabledPolicyIPListIdWithType(this.AdminContext(), params.FirewallPolicyId, params.Type)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Data["listId"] = listId

	// 数量
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
