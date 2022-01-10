package waf

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/dao"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs/firewallconfigs"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/maps"
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
		var actionMaps = []maps.Map{}
		for _, action := range set.Actions {
			def := firewallconfigs.FindActionDefinition(action.Code)
			if def == nil {
				continue
			}

			actionMaps = append(actionMaps, maps.Map{
				"code":     strings.ToUpper(action.Code),
				"name":     def.Name,
				"category": def.Category,
				"options":  action.Options,
			})
		}

		return maps.Map{
			"id":   set.Id,
			"name": set.Name,
			"rules": lists.Map(set.Rules, func(k int, v interface{}) interface{} {
				rule := v.(*firewallconfigs.HTTPFirewallRule)

				// 校验
				var errString = ""
				var err = rule.Init()
				if err != nil {
					errString = err.Error()
				}

				return maps.Map{
					"param":             rule.Param,
					"paramFilters":      rule.ParamFilters,
					"operator":          rule.Operator,
					"value":             rule.Value,
					"isCaseInsensitive": rule.IsCaseInsensitive,
					"description":       rule.Description,
					"isComposed":        firewallconfigs.CheckCheckpointIsComposed(rule.Prefix()),
					"checkpointOptions": rule.CheckpointOptions,
					"err":               errString,
				}
			}),
			"isOn":      set.IsOn,
			"actions":   actionMaps,
			"connector": strings.ToUpper(set.Connector),
		}
	})

	this.Show()
}
