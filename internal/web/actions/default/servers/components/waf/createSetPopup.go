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
	"strconv"
)

type CreateSetPopupAction struct {
	actionutils.ParentAction
}

func (this *CreateSetPopupAction) Init() {
	this.Nav("", "", "")
}

func (this *CreateSetPopupAction) RunGet(params struct {
	FirewallPolicyId int64
	GroupId          int64
	Type             string
}) {
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

	// 所有可选的动作
	var actionMaps = []maps.Map{}
	for _, action := range firewallconfigs.AllActions {
		actionMaps = append(actionMaps, maps.Map{
			"name":        action.Name,
			"description": action.Description,
			"code":        action.Code,
		})
	}
	this.Data["actions"] = actionMaps

	// 是否为全局
	this.Data["isGlobalPolicy"] = firewallPolicy.ServerId == 0

	this.Show()
}

func (this *CreateSetPopupAction) RunPost(params struct {
	GroupId int64

	Name string

	FormType string

	// normal
	RulesJSON          []byte
	Connector          string
	ActionsJSON        []byte
	IgnoreLocal        bool
	IgnoreSearchEngine bool

	// code
	Code string

	Must *actions.Must
}) {
	groupConfig, err := dao.SharedHTTPFirewallRuleGroupDAO.FindRuleGroupConfig(this.AdminContext(), params.GroupId)
	if err != nil {
		this.ErrorPage(err)
		return
	}
	if groupConfig == nil {
		this.Fail("找不到分组，Id：" + strconv.FormatInt(params.GroupId, 10))
		return
	}

	params.Must.
		Field("name", params.Name).
		Require("请输入规则集名称")

	var setConfigJSON []byte
	if params.FormType == "normal" {
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

		var setConfig = &firewallconfigs.HTTPFirewallRuleSet{
			Id:                 0,
			IsOn:               true,
			Name:               params.Name,
			Code:               "",
			Description:        "",
			Connector:          params.Connector,
			RuleRefs:           nil,
			Rules:              rules,
			Actions:            actionConfigs,
			IgnoreLocal:        params.IgnoreLocal,
			IgnoreSearchEngine: params.IgnoreSearchEngine,
		}

		setConfigJSON, err = json.Marshal(setConfig)
		if err != nil {
			this.ErrorPage(err)
			return
		}
	} else if params.FormType == "code" {
		var codeJSON = []byte(params.Code)
		if len(codeJSON) == 0 {
			this.FailField("code", "请输入规则集代码")
			return
		}

		var setConfig = &firewallconfigs.HTTPFirewallRuleSet{}
		err = json.Unmarshal(codeJSON, setConfig)
		if err != nil {
			this.FailField("code", "解析规则集代码失败："+err.Error())
			return
		}

		if len(setConfig.Rules) == 0 {
			this.FailField("code", "规则集代码中必须包含至少一个规则")
			return
		}

		if len(setConfig.Actions) == 0 {
			this.FailField("code", "规则集代码中必须包含至少一个动作")
			return
		}

		setConfig.Name = params.Name
		setConfig.IsOn = true

		// 重置ID
		setConfig.Id = 0

		setConfig.RuleRefs = nil
		for _, rule := range setConfig.Rules {
			rule.Id = 0
		}

		err = setConfig.Init()
		if err != nil {
			this.FailField("code", "校验规则集代码失败："+err.Error())
			return
		}

		setConfigJSON, err = json.Marshal(setConfig)
	} else {
		this.Fail("错误的参数'formType': " + params.FormType)
		return
	}

	createUpdateResp, err := this.RPC().HTTPFirewallRuleSetRPC().CreateOrUpdateHTTPFirewallRuleSetFromConfig(this.AdminContext(), &pb.CreateOrUpdateHTTPFirewallRuleSetFromConfigRequest{FirewallRuleSetConfigJSON: setConfigJSON})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	groupConfig.SetRefs = append(groupConfig.SetRefs, &firewallconfigs.HTTPFirewallRuleSetRef{
		IsOn:  true,
		SetId: createUpdateResp.FirewallRuleSetId,
	})

	setRefsJSON, err := json.Marshal(groupConfig.SetRefs)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	_, err = this.RPC().HTTPFirewallRuleGroupRPC().UpdateHTTPFirewallRuleGroupSets(this.AdminContext(), &pb.UpdateHTTPFirewallRuleGroupSetsRequest{
		FirewallRuleGroupId:  params.GroupId,
		FirewallRuleSetsJSON: setRefsJSON,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Data["setId"] = createUpdateResp.FirewallRuleSetId

	this.Success()
}
