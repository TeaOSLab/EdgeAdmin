package waf

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/dao"
	"github.com/iwind/TeaGo/maps"
)

type GroupsAction struct {
	actionutils.ParentAction
}

func (this *GroupsAction) Init() {
	this.Nav("", "", this.ParamString("type"))
}

func (this *GroupsAction) RunGet(params struct {
	FirewallPolicyId int64
	Type             string
}) {
	this.Data["type"] = params.Type

	firewallPolicy, err := dao.SharedHTTPFirewallPolicyDAO.FindEnabledHTTPFirewallPolicyConfig(this.AdminContext(), params.FirewallPolicyId)
	if err != nil {
		this.ErrorPage(err)
		return
	}
	if firewallPolicy == nil {
		this.NotFound("firewallPolicy", params.FirewallPolicyId)
		return
	}

	groupMaps := []maps.Map{}

	// inbound
	if params.Type == "inbound" {
		if firewallPolicy.Inbound != nil {
			for _, g := range firewallPolicy.Inbound.Groups {
				groupMaps = append(groupMaps, maps.Map{
					"id":          g.Id,
					"name":        g.Name,
					"code":        g.Code,
					"isOn":        g.IsOn,
					"description": g.Description,
					"countSets":   len(g.Sets),
					"isTemplate":  g.IsTemplate,
					"canDelete":   !g.IsTemplate,
				})
			}
		}
	}

	// outbound
	if params.Type == "outbound" {
		if firewallPolicy.Outbound != nil {
			for _, g := range firewallPolicy.Outbound.Groups {
				groupMaps = append(groupMaps, maps.Map{
					"id":          g.Id,
					"name":        g.Name,
					"code":        g.Code,
					"isOn":        g.IsOn,
					"description": g.Description,
					"countSets":   len(g.Sets),
					"isTemplate":  g.IsTemplate,
					"canDelete":   !g.IsTemplate,
				})
			}
		}
	}

	this.Data["groups"] = groupMaps

	this.Show()
}
