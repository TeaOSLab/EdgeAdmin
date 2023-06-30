package waf

import (	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/dao"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs/firewallconfigs"
	"github.com/iwind/TeaGo/actions"
)

type CreateGroupPopupAction struct {
	actionutils.ParentAction
}

func (this *CreateGroupPopupAction) Init() {
	this.Nav("", "", "")
}

func (this *CreateGroupPopupAction) RunGet(params struct {
	Type string
}) {
	this.Data["type"] = params.Type

	this.Show()
}

func (this *CreateGroupPopupAction) RunPost(params struct {
	FirewallPolicyId int64
	Type             string

	Name        string
	Code        string
	Description string
	IsOn        bool

	Must *actions.Must
}) {
	firewallPolicy, err := dao.SharedHTTPFirewallPolicyDAO.FindEnabledHTTPFirewallPolicyConfig(this.AdminContext(), params.FirewallPolicyId)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	if firewallPolicy == nil {
		this.NotFound("firewallPolicy", params.FirewallPolicyId)
	}

	params.Must.
		Field("name", params.Name).
		Require("请输入分组名称")

	createResp, err := this.RPC().HTTPFirewallRuleGroupRPC().CreateHTTPFirewallRuleGroup(this.AdminContext(), &pb.CreateHTTPFirewallRuleGroupRequest{
		IsOn:        params.IsOn,
		Name:        params.Name,
		Code:        params.Code,
		Description: params.Description,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	groupId := createResp.FirewallRuleGroupId

	switch params.Type {
	case "inbound":
		firewallPolicy.Inbound.GroupRefs = append(firewallPolicy.Inbound.GroupRefs, &firewallconfigs.HTTPFirewallRuleGroupRef{
			IsOn:    true,
			GroupId: groupId,
		})
	default:
		firewallPolicy.Outbound.GroupRefs = append(firewallPolicy.Outbound.GroupRefs, &firewallconfigs.HTTPFirewallRuleGroupRef{
			IsOn:    true,
			GroupId: groupId,
		})
	}

	inboundJSON, err := firewallPolicy.InboundJSON()
	if err != nil {
		this.ErrorPage(err)
		return
	}

	outboundJSON, err := firewallPolicy.OutboundJSON()
	if err != nil {
		this.ErrorPage(err)
		return
	}

	_, err = this.RPC().HTTPFirewallPolicyRPC().UpdateHTTPFirewallPolicyGroups(this.AdminContext(), &pb.UpdateHTTPFirewallPolicyGroupsRequest{
		HttpFirewallPolicyId: params.FirewallPolicyId,
		InboundJSON:          inboundJSON,
		OutboundJSON:         outboundJSON,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	// 日志
	defer this.CreateLogInfo(codes.WAFRuleGroup_LogCreateRuleGroup, groupId, params.Name)

	this.Success()
}
