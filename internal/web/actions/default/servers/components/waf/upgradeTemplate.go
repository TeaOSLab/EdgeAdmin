// Copyright 2021 GoEdge CDN goedge.cdn@gmail.com. All rights reserved.

package waf

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/dao"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs/firewallconfigs"
)

type UpgradeTemplateAction struct {
	actionutils.ParentAction
}

func (this *UpgradeTemplateAction) RunPost(params struct {
	PolicyId int64
}) {
	defer this.CreateLogInfo(codes.WAFPolicy_LogUpgradeWAFPolicy, params.PolicyId)

	policy, err := dao.SharedHTTPFirewallPolicyDAO.FindEnabledHTTPFirewallPolicyConfig(this.AdminContext(), params.PolicyId)
	if err != nil {
		this.ErrorPage(err)
		return
	}
	if policy == nil {
		this.NotFound("firewallPolicy", params.PolicyId)
		return
	}

	// 检查是否有升级
	var templatePolicy = firewallconfigs.HTTPFirewallTemplate()
	if templatePolicy.Inbound != nil {
		for _, group := range templatePolicy.Inbound.Groups {
			if len(group.Code) == 0 {
				continue
			}
			var oldGroup = policy.FindRuleGroupWithCode(group.Code)
			if oldGroup == nil {
				createGroupResp, err := this.RPC().HTTPFirewallRuleGroupRPC().CreateHTTPFirewallRuleGroup(this.AdminContext(), &pb.CreateHTTPFirewallRuleGroupRequest{
					IsOn:        true,
					Name:        group.Name,
					Code:        group.Code,
					Description: group.Description,
				})
				if err != nil {
					this.ErrorPage(err)
					return
				}
				var groupId = createGroupResp.FirewallRuleGroupId
				policy.Inbound.GroupRefs = append(policy.Inbound.GroupRefs, &firewallconfigs.HTTPFirewallRuleGroupRef{
					IsOn:    true,
					GroupId: groupId,
				})

				for _, set := range group.Sets {
					setJSON, err := json.Marshal(set)
					if err != nil {
						this.ErrorPage(err)
						return
					}
					_, err = this.RPC().HTTPFirewallRuleGroupRPC().AddHTTPFirewallRuleGroupSet(this.AdminContext(), &pb.AddHTTPFirewallRuleGroupSetRequest{
						FirewallRuleGroupId:       groupId,
						FirewallRuleSetConfigJSON: setJSON,
					})
					if err != nil {
						this.ErrorPage(err)
						return
					}
				}

				continue
			}
			for _, set := range group.Sets {
				if len(set.Code) == 0 {
					continue
				}
				var oldSet = oldGroup.FindRuleSetWithCode(set.Code)
				if oldSet == nil {
					setJSON, err := json.Marshal(set)
					if err != nil {
						this.ErrorPage(err)
						return
					}
					_, err = this.RPC().HTTPFirewallRuleGroupRPC().AddHTTPFirewallRuleGroupSet(this.AdminContext(), &pb.AddHTTPFirewallRuleGroupSetRequest{
						FirewallRuleGroupId:       oldGroup.Id,
						FirewallRuleSetConfigJSON: setJSON,
					})
					if err != nil {
						this.ErrorPage(err)
						return
					}
					continue
				}
			}
		}
	}

	// 保存inbound
	inboundJSON, err := policy.InboundJSON()
	if err != nil {
		this.ErrorPage(err)
		return
	}

	outboundJSON, err := policy.OutboundJSON()
	if err != nil {
		this.ErrorPage(err)
		return
	}

	_, err = this.RPC().HTTPFirewallPolicyRPC().UpdateHTTPFirewallPolicyGroups(this.AdminContext(), &pb.UpdateHTTPFirewallPolicyGroupsRequest{
		HttpFirewallPolicyId: params.PolicyId,
		InboundJSON:          inboundJSON,
		OutboundJSON:         outboundJSON,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
