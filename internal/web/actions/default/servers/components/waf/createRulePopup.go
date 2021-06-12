package waf

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs/firewallconfigs"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/maps"
)

type CreateRulePopupAction struct {
	actionutils.ParentAction
}

func (this *CreateRulePopupAction) Init() {
	this.Nav("", "", "")
}

func (this *CreateRulePopupAction) RunGet(params struct {
	Type string
}) {
	// check points
	checkpointList := []maps.Map{}
	for _, checkpoint := range firewallconfigs.AllCheckpoints {
		if (params.Type == "inbound" && checkpoint.IsRequest) || params.Type == "outbound" {
			checkpointList = append(checkpointList, maps.Map{
				"name":        checkpoint.Name,
				"prefix":      checkpoint.Prefix,
				"description": checkpoint.Description,
				"hasParams":   checkpoint.HasParams,
				"params":      checkpoint.Params,
				"options":     checkpoint.Options,
				"isComposed":  checkpoint.IsComposed,
			})
		}
	}

	// operators
	this.Data["operators"] = lists.Map(firewallconfigs.AllRuleOperators, func(k int, v interface{}) interface{} {
		def := v.(*firewallconfigs.RuleOperatorDefinition)
		return maps.Map{
			"name":        def.Name,
			"code":        def.Code,
			"description": def.Description,
			"case":        def.CaseInsensitive,
		}
	})

	this.Data["checkpoints"] = checkpointList

	this.Show()
}

func (this *CreateRulePopupAction) RunPost(params struct {
	RuleId           int64
	Prefix           string
	Operator         string
	Param            string
	ParamFiltersJSON []byte
	OptionsJSON      []byte
	Value            string
	Case             bool

	Must *actions.Must
}) {
	params.Must.
		Field("prefix", params.Prefix).
		Require("请选择参数")

	rule := &firewallconfigs.HTTPFirewallRule{
		Id:   params.RuleId,
		IsOn: true,
	}
	if len(params.Param) > 0 {
		rule.Param = "${" + params.Prefix + "." + params.Param + "}"
	} else {
		rule.Param = "${" + params.Prefix + "}"
	}

	paramFilters := []*firewallconfigs.ParamFilter{}
	if len(params.ParamFiltersJSON) > 0 {
		err := json.Unmarshal(params.ParamFiltersJSON, &paramFilters)
		if err != nil {
			this.ErrorPage(err)
			return
		}
	}
	rule.ParamFilters = paramFilters

	rule.Operator = params.Operator
	rule.Value = params.Value
	rule.IsCaseInsensitive = params.Case

	if len(params.OptionsJSON) > 0 {
		options := []maps.Map{}
		err := json.Unmarshal(params.OptionsJSON, &options)
		if err != nil {
			this.ErrorPage(err)
			return
		}

		rule.CheckpointOptions = map[string]interface{}{}
		for _, option := range options {
			rule.CheckpointOptions[option.GetString("code")] = option.Get("value")
		}
	}

	// 校验
	err := rule.Init()
	if err != nil {
		this.Fail("校验规则 '" + rule.Param + " " + rule.Operator + " " + rule.Value + "' 失败，原因：" + err.Error())
	}

	this.Data["rule"] = rule
	this.Success()
}
