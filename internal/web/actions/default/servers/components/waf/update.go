package waf

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/models"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs/firewallconfigs"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type UpdateAction struct {
	actionutils.ParentAction
}

func (this *UpdateAction) Init() {
	this.Nav("", "", "update")
}

func (this *UpdateAction) RunGet(params struct {
	FirewallPolicyId int64
}) {
	firewallPolicy, err := models.SharedHTTPFirewallPolicyDAO.FindEnabledPolicyConfig(this.AdminContext(), params.FirewallPolicyId)
	if err != nil {
		this.ErrorPage(err)
		return
	}
	if firewallPolicy == nil {
		this.NotFound("firewallPolicy", params.FirewallPolicyId)
		return
	}
	this.Data["firewallPolicy"] = maps.Map{
		"id":          firewallPolicy.Id,
		"name":        firewallPolicy.Name,
		"description": firewallPolicy.Description,
		"isOn":        firewallPolicy.IsOn,
	}

	// 预置分组
	groups := []maps.Map{}
	templatePolicy := firewallconfigs.HTTPFirewallTemplate()
	for _, group := range templatePolicy.AllRuleGroups() {
		if len(group.Code) > 0 {
			usedGroup := firewallPolicy.FindRuleGroupWithCode(group.Code)
			if usedGroup != nil {
				group.IsOn = usedGroup.IsOn
			}
		}

		groups = append(groups, maps.Map{
			"code": group.Code,
			"name": group.Name,
			"isOn": group.IsOn,
		})
	}
	this.Data["groups"] = groups

	this.Show()
}

func (this *UpdateAction) RunPost(params struct {
	FirewallPolicyId int64
	Name             string
	GroupCodes       []string
	Description      string
	IsOn             bool

	Must *actions.Must
}) {
	params.Must.
		Field("name", params.Name).
		Require("请输入策略名称")

	_, err := this.RPC().HTTPFirewallPolicyRPC().UpdateHTTPFirewallPolicy(this.AdminContext(), &pb.UpdateHTTPFirewallPolicyRequest{
		FirewallPolicyId:   params.FirewallPolicyId,
		IsOn:               params.IsOn,
		Name:               params.Name,
		Description:        params.Description,
		FirewallGroupCodes: params.GroupCodes,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
