package waf

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs/firewallconfigs"
	"github.com/iwind/TeaGo/maps"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.FirstMenu("index")
}

func (this *IndexAction) RunGet(params struct{}) {
	countResp, err := this.RPC().HTTPFirewallPolicyRPC().CountAllEnabledFirewallPolicies(this.AdminContext(), &pb.CountAllEnabledFirewallPoliciesRequest{})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	count := countResp.Count
	page := this.NewPage(count)

	listResp, err := this.RPC().HTTPFirewallPolicyRPC().ListEnabledFirewallPolicies(this.AdminContext(), &pb.ListEnabledFirewallPoliciesRequest{
		Offset: page.Offset,
		Size:   page.Size,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	policyMaps := []maps.Map{}
	for _, policy := range listResp.FirewallPolicies {
		countInbound := 0
		countOutbound := 0
		if len(policy.InboundJSON) > 0 {
			inboundConfig := &firewallconfigs.HTTPFirewallInboundConfig{}
			err = json.Unmarshal(policy.InboundJSON, inboundConfig)
			if err != nil {
				this.ErrorPage(err)
				return
			}
			countInbound = len(inboundConfig.GroupRefs)
		}
		if len(policy.OutboundJSON) > 0 {
			outboundConfig := &firewallconfigs.HTTPFirewallInboundConfig{}
			err = json.Unmarshal(policy.OutboundJSON, outboundConfig)
			if err != nil {
				this.ErrorPage(err)
				return
			}
			countOutbound = len(outboundConfig.GroupRefs)
		}

		countServersResp, err := this.RPC().ServerRPC().CountAllEnabledServersWithHTTPFirewallPolicyId(this.AdminContext(), &pb.CountAllEnabledServersWithHTTPFirewallPolicyIdRequest{FirewallPolicyId: policy.Id})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		countServers := countServersResp.Count

		policyMaps = append(policyMaps, maps.Map{
			"id":            policy.Id,
			"isOn":          policy.IsOn,
			"name":          policy.Name,
			"countInbound":  countInbound,
			"countOutbound": countOutbound,
			"countServers":  countServers,
		})
	}

	this.Data["policies"] = policyMaps

	this.Data["page"] = page.AsHTML()

	this.Show()
}
