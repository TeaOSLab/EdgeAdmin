package waf

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/models"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs/firewallconfigs"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
	"strconv"
	"strings"
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

	firewallPolicy, err := models.SharedHTTPFirewallPolicyDAO.FindEnabledPolicyConfig(this.AdminContext(), params.FirewallPolicyId)
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

	this.Show()
}

func (this *CreateSetPopupAction) RunPost(params struct {
	GroupId int64

	Name      string
	RulesJSON []byte
	Connector string
	Action    string

	Must *actions.Must
}) {
	groupConfig, err := models.SharedHTTPFirewallRuleGroupDAO.FindRuleGroupConfig(this.AdminContext(), params.GroupId)
	if err != nil {
		this.ErrorPage(err)
		return
	}
	if groupConfig == nil {
		this.Fail("找不到分组，Id：" + strconv.FormatInt(params.GroupId, 10))
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
	}
	if len(rules) == 0 {
		this.Fail("请添加至少一个规则")
	}

	setConfig := &firewallconfigs.HTTPFirewallRuleSet{
		Id:            0,
		IsOn:          true,
		Name:          params.Name,
		Code:          "",
		Description:   "",
		Connector:     params.Connector,
		RuleRefs:      nil,
		Rules:         rules,
		Action:        params.Action,
		ActionOptions: maps.Map{},
	}

	for k, v := range this.ParamsMap {
		if len(v) == 0 {
			continue
		}
		index := strings.Index(k, "action_")
		if index > -1 {
			setConfig.ActionOptions[k[len("action_"):]] = v[0]
		}
	}

	setConfigJSON, err := json.Marshal(setConfig)
	if err != nil {
		this.ErrorPage(err)
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

	this.Success()
}
