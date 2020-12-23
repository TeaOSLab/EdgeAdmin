package waf

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/dao"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs/firewallconfigs"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/maps"
	"strconv"
	"strings"
)

type GroupAction struct {
	actionutils.ParentAction
}

func (this *GroupAction) Init() {
	this.Nav("", "", this.ParamString("type"))
}

func (this *GroupAction) RunGet(params struct {
	FirewallPolicyId int64
	GroupId          int64
	Type             string
}) {
	this.Data["type"] = params.Type

	// policy
	firewallPolicy, err := dao.SharedHTTPFirewallPolicyDAO.FindEnabledHTTPFirewallPolicyConfig(this.AdminContext(), params.FirewallPolicyId)
	if err != nil {
		this.ErrorPage(err)
		return
	}
	if firewallPolicy == nil {
		this.NotFound("firewallPolicy", params.FirewallPolicyId)
		return
	}

	// group config
	groupConfig, err := dao.SharedHTTPFirewallRuleGroupDAO.FindRuleGroupConfig(this.AdminContext(), params.GroupId)
	if err != nil {
		this.ErrorPage(err)
		return
	}
	if groupConfig == nil {
		this.NotFound("firewallRuleGroup", params.GroupId)
		return
	}

	this.Data["group"] = groupConfig

	// rule sets
	this.Data["sets"] = lists.Map(groupConfig.Sets, func(k int, v interface{}) interface{} {
		set := v.(*firewallconfigs.HTTPFirewallRuleSet)

		// 动作说明
		actionLinks := []maps.Map{}
		if set.Action == firewallconfigs.HTTPFirewallActionGoGroup {
			nextGroup := firewallPolicy.FindRuleGroup(set.ActionOptions.GetInt64("groupId"))
			if nextGroup != nil {
				actionLinks = append(actionLinks, maps.Map{
					"name": nextGroup.Name,
					"url":  "/servers/components/waf/group?firewallPolicyId=" + strconv.FormatInt(params.FirewallPolicyId, 10) + "&type=" + params.Type + "&groupId=" + strconv.FormatInt(nextGroup.Id, 10),
				})
			}
		} else if set.Action == firewallconfigs.HTTPFirewallActionGoSet {
			nextGroup := firewallPolicy.FindRuleGroup(set.ActionOptions.GetInt64("groupId"))
			if nextGroup != nil {
				actionLinks = append(actionLinks, maps.Map{
					"name": nextGroup.Name,
					"url":  "/servers/components/waf/group?firewallPolicyId=" + strconv.FormatInt(params.FirewallPolicyId, 10) + "&type=" + params.Type + "&groupId=" + strconv.FormatInt(nextGroup.Id, 10),
				})

				nextSet := nextGroup.FindRuleSet(set.ActionOptions.GetInt64("setId"))
				if nextSet != nil {
					actionLinks = append(actionLinks, maps.Map{
						"name": nextSet.Name,
						"url":  "/servers/components/waf/group?firewallPolicyId=" + strconv.FormatInt(params.FirewallPolicyId, 10) + "&type=" + params.Type + "&groupId=" + strconv.FormatInt(nextGroup.Id, 10),
					})
				}
			}
		}

		return maps.Map{
			"id":   set.Id,
			"name": set.Name,
			"rules": lists.Map(set.Rules, func(k int, v interface{}) interface{} {
				rule := v.(*firewallconfigs.HTTPFirewallRule)
				return maps.Map{
					"param":             rule.Param,
					"paramFilters":      rule.ParamFilters,
					"operator":          rule.Operator,
					"value":             rule.Value,
					"isCaseInsensitive": rule.IsCaseInsensitive,
					"isComposed":        firewallconfigs.CheckCheckpointIsComposed(rule.Prefix()),
				}
			}),
			"isOn":          set.IsOn,
			"action":        strings.ToUpper(set.Action),
			"actionOptions": set.ActionOptions,
			"actionName":    firewallconfigs.FindActionName(set.Action),
			"actionLinks":   actionLinks,
			"connector":     strings.ToUpper(set.Connector),
		}
	})

	this.Show()
}
