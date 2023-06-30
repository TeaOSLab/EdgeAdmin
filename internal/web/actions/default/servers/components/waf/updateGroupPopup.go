package waf

import (	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/dao"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type UpdateGroupPopupAction struct {
	actionutils.ParentAction
}

func (this *UpdateGroupPopupAction) Init() {
	this.Nav("", "", "")
}

func (this *UpdateGroupPopupAction) RunGet(params struct {
	GroupId int64
}) {
	groupConfig, err := dao.SharedHTTPFirewallRuleGroupDAO.FindRuleGroupConfig(this.AdminContext(), params.GroupId)
	if err != nil {
		this.ErrorPage(err)
		return
	}
	if groupConfig == nil {
		this.NotFound("ruleGroup", params.GroupId)
		return
	}

	this.Data["group"] = maps.Map{
		"id":          groupConfig.Id,
		"name":        groupConfig.Name,
		"description": groupConfig.Description,
		"isOn":        groupConfig.IsOn,
		"code":        groupConfig.Code,
	}

	this.Show()
}

func (this *UpdateGroupPopupAction) RunPost(params struct {
	GroupId     int64
	Name        string
	Code        string
	Description string
	IsOn        bool

	Must *actions.Must
}) {
	// 日志
	defer this.CreateLogInfo(codes.WAFRuleGroup_LogUpdateRuleGroup, params.GroupId)

	params.Must.
		Field("name", params.Name).
		Require("请输入分组名称")

	_, err := this.RPC().HTTPFirewallRuleGroupRPC().UpdateHTTPFirewallRuleGroup(this.AdminContext(), &pb.UpdateHTTPFirewallRuleGroupRequest{
		FirewallRuleGroupId: params.GroupId,
		IsOn:                params.IsOn,
		Name:                params.Name,
		Code:                params.Code,
		Description:         params.Description,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
