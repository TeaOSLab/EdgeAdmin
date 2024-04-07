package waf

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/utils/numberutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/dao"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs/firewallconfigs"
	"github.com/iwind/TeaGo/maps"
)

type PolicyAction struct {
	actionutils.ParentAction
}

func (this *PolicyAction) Init() {
	this.Nav("", "", "index")
}

func (this *PolicyAction) RunGet(params struct {
	FirewallPolicyId int64
}) {
	firewallPolicy, err := dao.SharedHTTPFirewallPolicyDAO.FindEnabledHTTPFirewallPolicyConfig(this.AdminContext(), params.FirewallPolicyId)
	if err != nil {
		this.ErrorPage(err)
		return
	}
	if firewallPolicy == nil {
		this.NotFound("firewallPolicy", params.FirewallPolicyId)
		return
	}

	var internalGroups = []maps.Map{}
	if firewallPolicy.Inbound != nil {
		for _, group := range firewallPolicy.Inbound.Groups {
			internalGroups = append(internalGroups, maps.Map{
				"name": group.Name,
				"isOn": group.IsOn,
			})
		}
	}
	if firewallPolicy.Outbound != nil {
		for _, group := range firewallPolicy.Outbound.Groups {
			internalGroups = append(internalGroups, maps.Map{
				"name": group.Name,
				"isOn": group.IsOn,
			})
		}
	}

	// 检查是否有升级
	var templatePolicy = firewallconfigs.HTTPFirewallTemplate()
	var upgradeItems = []maps.Map{}
	if templatePolicy.Inbound != nil {
		for _, group := range templatePolicy.Inbound.Groups {
			if len(group.Code) == 0 {
				continue
			}
			var oldGroup = firewallPolicy.FindRuleGroupWithCode(group.Code)
			if oldGroup == nil {
				upgradeItems = append(upgradeItems, maps.Map{
					"name": group.Name,
					"isOn": group.IsOn,
				})
				continue
			}
			for _, set := range group.Sets {
				if len(set.Code) == 0 {
					continue
				}
				var oldSet = oldGroup.FindRuleSetWithCode(set.Code)
				if oldSet == nil {
					upgradeItems = append(upgradeItems, maps.Map{
						"name": group.Name + " -- " + set.Name,
						"isOn": set.IsOn,
					})
					continue
				}
			}
		}
	}
	this.Data["upgradeItems"] = upgradeItems

	// 模式
	if len(firewallPolicy.Mode) == 0 {
		firewallPolicy.Mode = firewallconfigs.FirewallModeDefend
	}

	this.Data["firewallPolicy"] = maps.Map{
		"id":                       firewallPolicy.Id,
		"name":                     firewallPolicy.Name,
		"isOn":                     firewallPolicy.IsOn,
		"description":              firewallPolicy.Description,
		"mode":                     firewallPolicy.Mode,
		"modeInfo":                 firewallconfigs.FindFirewallMode(firewallPolicy.Mode),
		"groups":                   internalGroups,
		"blockOptions":             firewallPolicy.BlockOptions,
		"pageOptions":              firewallPolicy.PageOptions,
		"captchaOptions":           firewallPolicy.CaptchaOptions,
		"jsCookieOptions":          firewallPolicy.JSCookieOptions,
		"useLocalFirewall":         firewallPolicy.UseLocalFirewall,
		"synFlood":                 firewallPolicy.SYNFlood,
		"log":                      firewallPolicy.Log,
		"maxRequestBodySize":       firewallPolicy.MaxRequestBodySize,
		"maxRequestBodySizeFormat": numberutils.FormatBytes(firewallPolicy.MaxRequestBodySize),
		"denyCountryHTML":          firewallPolicy.DenyCountryHTML,
		"denyProvinceHTML":         firewallPolicy.DenyProvinceHTML,
	}

	// 正在使用此策略的集群
	clustersResp, err := this.RPC().NodeClusterRPC().FindAllEnabledNodeClustersWithHTTPFirewallPolicyId(this.AdminContext(), &pb.FindAllEnabledNodeClustersWithHTTPFirewallPolicyIdRequest{HttpFirewallPolicyId: params.FirewallPolicyId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	clusterMaps := []maps.Map{}
	for _, cluster := range clustersResp.NodeClusters {
		clusterMaps = append(clusterMaps, maps.Map{
			"id":   cluster.Id,
			"name": cluster.Name,
		})
	}
	this.Data["clusters"] = clusterMaps

	this.Show()
}
