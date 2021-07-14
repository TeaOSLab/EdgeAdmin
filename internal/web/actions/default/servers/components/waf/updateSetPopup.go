package waf

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/oplogs"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/dao"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs/firewallconfigs"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type UpdateSetPopupAction struct {
	actionutils.ParentAction
}

func (this *UpdateSetPopupAction) Init() {
	this.Nav("", "", "")
}

func (this *UpdateSetPopupAction) RunGet(params struct {
	FirewallPolicyId int64
	GroupId          int64
	Type             string
	SetId            int64
}) {
	// 日志
	defer this.CreateLog(oplogs.LevelInfo, "修改WAF规则集 %d 基本信息", params.SetId)

	this.Data["groupId"] = params.GroupId
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
	this.Data["firewallPolicy"] = firewallPolicy

	// 一些配置
	this.Data["connectors"] = []maps.Map{
		{
			"name":        "和(AND)",
			"value":       firewallconfigs.HTTPFirewallRuleConnectorAnd,
			"description": "所有规则都满足才视为匹配",
		},
		{
			"name":        "或(OR)",
			"value":       firewallconfigs.HTTPFirewallRuleConnectorOr,
			"description": "任一规则满足了就视为匹配",
		},
	}

	actionMaps := []maps.Map{}
	for _, action := range firewallconfigs.AllActions {
		actionMaps = append(actionMaps, maps.Map{
			"name":        action.Name,
			"description": action.Description,
			"code":        action.Code,
		})
	}
	this.Data["actions"] = actionMaps

	// 规则集信息
	setConfig, err := dao.SharedHTTPFirewallRuleSetDAO.FindRuleSetConfig(this.AdminContext(), params.SetId)
	if err != nil {
		this.ErrorPage(err)
		return
	}
	if setConfig == nil {
		this.NotFound("firewallRuleSet", params.SetId)
		return
	}
	this.Data["setConfig"] = setConfig

	// action configs
	actionConfigs, err := dao.SharedHTTPFirewallPolicyDAO.FindHTTPFirewallActionConfigs(this.AdminContext(), setConfig.Actions)
	if err != nil {
		this.ErrorPage(err)
		return
	}
	this.Data["actionConfigs"] = actionConfigs

	this.Show()
}

func (this *UpdateSetPopupAction) RunPost(params struct {
	GroupId int64
	SetId   int64

	Name        string
	RulesJSON   []byte
	Connector   string
	ActionsJSON []byte

	Must *actions.Must
}) {
	// 规则集信息
	setConfig, err := dao.SharedHTTPFirewallRuleSetDAO.FindRuleSetConfig(this.AdminContext(), params.SetId)
	if err != nil {
		this.ErrorPage(err)
		return
	}
	if setConfig == nil {
		this.NotFound("firewallRuleSet", params.SetId)
		return
	}

	params.Must.
		Field("name", params.Name).
		Require("请输入规则集名称")

	if len(params.RulesJSON) == 0 {
		this.Fail("请添加至少一个规则")
	}
	rules := []*firewallconfigs.HTTPFirewallRule{}
	err = json.Unmarshal(params.RulesJSON, &rules)
	if err != nil {
		this.ErrorPage(err)
		return
	}
	if len(rules) == 0 {
		this.Fail("请添加至少一个规则")
	}

	var actionConfigs = []*firewallconfigs.HTTPFirewallActionConfig{}
	if len(params.ActionsJSON) > 0 {
		err = json.Unmarshal(params.ActionsJSON, &actionConfigs)
		if err != nil {
			this.ErrorPage(err)
			return
		}
	}
	if len(actionConfigs) == 0 {
		this.Fail("请添加至少一个动作")
	}

	setConfig.Name = params.Name
	setConfig.Connector = params.Connector
	setConfig.Rules = rules
	setConfig.Actions = actionConfigs

	setConfigJSON, err := json.Marshal(setConfig)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	_, err = this.RPC().HTTPFirewallRuleSetRPC().CreateOrUpdateHTTPFirewallRuleSetFromConfig(this.AdminContext(), &pb.CreateOrUpdateHTTPFirewallRuleSetFromConfigRequest{FirewallRuleSetConfigJSON: setConfigJSON})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
