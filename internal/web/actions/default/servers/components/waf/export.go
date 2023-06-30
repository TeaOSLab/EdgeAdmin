package waf

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/ttlcache"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/dao"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs/firewallconfigs"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/rands"
	"time"
)

type ExportAction struct {
	actionutils.ParentAction
}

func (this *ExportAction) Init() {
	this.Nav("", "", "export")
}

func (this *ExportAction) RunGet(params struct {
	FirewallPolicyId int64
}) {
	policy, err := dao.SharedHTTPFirewallPolicyDAO.FindEnabledHTTPFirewallPolicyConfig(this.AdminContext(), params.FirewallPolicyId)
	if err != nil {
		this.ErrorPage(err)
		return
	}
	if policy == nil {
		this.NotFound("firewallPolicy", policy.Id)
		return
	}

	enabledInboundGroups := []*firewallconfigs.HTTPFirewallRuleGroup{}
	enabledOutboundGroups := []*firewallconfigs.HTTPFirewallRuleGroup{}

	disabledInboundGroups := []*firewallconfigs.HTTPFirewallRuleGroup{}
	disabledOutboundGroups := []*firewallconfigs.HTTPFirewallRuleGroup{}

	if policy.Inbound != nil {
		for _, g := range policy.Inbound.Groups {
			if g.IsOn {
				enabledInboundGroups = append(enabledInboundGroups, g)
			} else {
				disabledInboundGroups = append(disabledInboundGroups, g)
			}
		}
	}
	if policy.Outbound != nil {
		for _, g := range policy.Outbound.Groups {
			if g.IsOn {
				enabledOutboundGroups = append(enabledOutboundGroups, g)
			} else {
				disabledOutboundGroups = append(disabledOutboundGroups, g)
			}
		}
	}

	this.Data["enabledInboundGroups"] = enabledInboundGroups
	this.Data["enabledOutboundGroups"] = enabledOutboundGroups

	this.Data["disabledInboundGroups"] = disabledInboundGroups
	this.Data["disabledOutboundGroups"] = disabledOutboundGroups

	this.Show()
}

func (this *ExportAction) RunPost(params struct {
	FirewallPolicyId int64
	InboundGroupIds  []int64
	OutboundGroupIds []int64

	Must *actions.Must
	CSRF *actionutils.CSRF
}) {
	defer this.CreateLogInfo(codes.WAFPolicy_LogExportWAFPolicy, params.FirewallPolicyId)

	policy, err := dao.SharedHTTPFirewallPolicyDAO.FindEnabledHTTPFirewallPolicyConfig(this.AdminContext(), params.FirewallPolicyId)
	if err != nil {
		this.ErrorPage(err)
		return
	}
	if policy == nil {
		this.NotFound("firewallPolicy", policy.Id)
		return
	}

	// inbound
	newInboundGroups := []*firewallconfigs.HTTPFirewallRuleGroup{}
	for _, inboundGroupId := range params.InboundGroupIds {
		group := policy.FindRuleGroup(inboundGroupId)
		if group != nil {
			newInboundGroups = append(newInboundGroups, group)
		}
	}
	if policy.Inbound == nil {
		policy.Inbound = &firewallconfigs.HTTPFirewallInboundConfig{
			IsOn: true,
		}
	}
	policy.Inbound.Groups = newInboundGroups
	policy.Inbound.GroupRefs = nil

	// outbound
	newOutboundGroups := []*firewallconfigs.HTTPFirewallRuleGroup{}
	for _, outboundGroupId := range params.OutboundGroupIds {
		group := policy.FindRuleGroup(outboundGroupId)
		if group != nil {
			newOutboundGroups = append(newOutboundGroups, group)
		}
	}
	if policy.Outbound == nil {
		policy.Outbound = &firewallconfigs.HTTPFirewallOutboundConfig{
			IsOn: true,
		}
	}
	policy.Outbound.Groups = newOutboundGroups
	policy.Outbound.GroupRefs = nil

	configJSON, err := json.Marshal(policy)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	key := "waf." + rands.HexString(32)
	ttlcache.DefaultCache.Write(key, configJSON, time.Now().Unix()+600)

	this.Data["key"] = key
	this.Data["id"] = params.FirewallPolicyId
	this.Success()
}
