package waf

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
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
	defer this.CreateLogInfo(codes.WAFRuleSet_LogUpdateRuleSet, params.SetId)

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
			"name":        this.Lang(codes.WAF_ConnectorAnd),
			"value":       firewallconfigs.HTTPFirewallRuleConnectorAnd,
			"description": this.Lang(codes.WAF_ConnectorAndDescription),
		},
		{
			"name":        this.Lang(codes.WAF_ConnectorOr),
			"value":       firewallconfigs.HTTPFirewallRuleConnectorOr,
			"description": this.Lang(codes.WAF_ConnectorOrDescription),
		},
	}

	var actionMaps = []maps.Map{}
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

	Name               string
	RulesJSON          []byte
	Connector          string
	ActionsJSON        []byte
	IgnoreLocal        bool
	IgnoreSearchEngine bool

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
		return
	}
	var rules = []*firewallconfigs.HTTPFirewallRule{}
	err = json.Unmarshal(params.RulesJSON, &rules)
	if err != nil {
		this.ErrorPage(err)
		return
	}
	if len(rules) == 0 {
		this.Fail("请添加至少一个规则")
		return
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
		return
	}

	setConfig.Name = params.Name
	setConfig.Connector = params.Connector
	setConfig.Rules = rules
	setConfig.Actions = actionConfigs
	setConfig.IgnoreLocal = params.IgnoreLocal
	setConfig.IgnoreSearchEngine = params.IgnoreSearchEngine

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

	this.Success()
}
