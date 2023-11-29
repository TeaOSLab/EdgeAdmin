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

func (this *IndexAction) RunGet(params struct {
	Keyword   string
	ClusterId int64
}) {
	this.Data["keyword"] = params.Keyword
	this.Data["clusterId"] = params.ClusterId

	countResp, err := this.RPC().HTTPFirewallPolicyRPC().CountAllEnabledHTTPFirewallPolicies(this.AdminContext(), &pb.CountAllEnabledHTTPFirewallPoliciesRequest{
		NodeClusterId: params.ClusterId,
		Keyword:       params.Keyword,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	var count = countResp.Count
	var page = this.NewPage(count)

	listResp, err := this.RPC().HTTPFirewallPolicyRPC().ListEnabledHTTPFirewallPolicies(this.AdminContext(), &pb.ListEnabledHTTPFirewallPoliciesRequest{
		NodeClusterId: params.ClusterId,
		Keyword:       params.Keyword,
		Offset:        page.Offset,
		Size:          page.Size,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	var policyMaps = []maps.Map{}
	for _, policy := range listResp.HttpFirewallPolicies {
		var countInbound = 0
		var countOutbound = 0
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

		countClustersResp, err := this.RPC().NodeClusterRPC().CountAllEnabledNodeClustersWithHTTPFirewallPolicyId(this.AdminContext(), &pb.CountAllEnabledNodeClustersWithHTTPFirewallPolicyIdRequest{HttpFirewallPolicyId: policy.Id})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		var countClusters = countClustersResp.Count

		// mode
		if len(policy.Mode) == 0 {
			policy.Mode = firewallconfigs.FirewallModeDefend
		}

		policyMaps = append(policyMaps, maps.Map{
			"id":            policy.Id,
			"isOn":          policy.IsOn,
			"name":          policy.Name,
			"mode":          policy.Mode,
			"modeInfo":      firewallconfigs.FindFirewallMode(policy.Mode),
			"countInbound":  countInbound,
			"countOutbound": countOutbound,
			"countClusters": countClusters,
		})
	}

	this.Data["policies"] = policyMaps

	this.Data["page"] = page.AsHTML()

	this.Show()
}
