package ipadmin

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/dao"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs/firewallconfigs"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs/ipconfigs"
	"github.com/iwind/TeaGo/maps"
	timeutil "github.com/iwind/TeaGo/utils/time"
	"time"
)

type GreyListAction struct {
	actionutils.ParentAction
}

func (this *GreyListAction) Init() {
	this.Nav("", "setting", "greyList")
	this.SecondMenu("waf")
}

func (this *GreyListAction) RunGet(params struct {
	FirewallPolicyId int64
	ServerId         int64
}) {
	this.Data["featureIsOn"] = true
	this.Data["firewallPolicyId"] = params.FirewallPolicyId

	listId, err := dao.SharedIPListDAO.FindGreyIPListIdWithServerId(this.AdminContext(), params.ServerId)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	// 创建
	if listId == 0 {
		listId, err = dao.SharedIPListDAO.CreateIPListForServerId(this.AdminContext(), params.ServerId, ipconfigs.IPListTypeGrey)
		if err != nil {
			this.ErrorPage(err)
			return
		}
	}

	this.Data["listId"] = listId

	// 数量
	countResp, err := this.RPC().IPItemRPC().CountIPItemsWithListId(this.AdminContext(), &pb.CountIPItemsWithListIdRequest{IpListId: listId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	var count = countResp.Count
	var page = this.NewPage(count)
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
	var itemMaps = []maps.Map{}
	for _, item := range itemsResp.IpItems {
		var expiredTime = ""
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
			"value":          item.Value,
			"ipFrom":         item.IpFrom,
			"ipTo":           item.IpTo,
			"createdTime":    timeutil.FormatTime("Y-m-d", item.CreatedAt),
			"expiredTime":    expiredTime,
			"lifeSeconds":    item.ExpiredAt - time.Now().Unix(),
			"reason":         item.Reason,
			"type":           item.Type,
			"isExpired":      item.ExpiredAt > 0 && item.ExpiredAt < time.Now().Unix(),
			"eventLevelName": firewallconfigs.FindFirewallEventLevelName(item.EventLevel),
			"sourcePolicy":   sourcePolicyMap,
			"sourceGroup":    sourceGroupMap,
			"sourceSet":      sourceSetMap,
			"sourceServer":   sourceServerMap,
		})
	}
	this.Data["items"] = itemMaps

	// WAF是否启用
	webConfig, err := dao.SharedHTTPWebDAO.FindWebConfigWithServerId(this.AdminContext(), params.ServerId)
	if err != nil {
		this.ErrorPage(err)
		return
	}
	this.Data["wafIsOn"] = webConfig.FirewallRef != nil && webConfig.FirewallRef.IsOn

	this.Show()
}
