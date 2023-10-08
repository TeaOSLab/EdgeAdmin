package waf

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs/firewallconfigs"
	"github.com/iwind/TeaGo/actions"
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
	var checkpointList = []maps.Map{}
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
				"dataType":    checkpoint.DataType,
			})
		}
	}

	// operators
	var operatorMaps = []maps.Map{}
	for _, operator := range firewallconfigs.AllRuleOperators {
		operatorMaps = append(operatorMaps, maps.Map{
			"name":        operator.Name,
			"code":        operator.Code,
			"description": operator.Description,
			"case":        operator.CaseInsensitive,
			"dataType":    operator.DataType,
		})
	}
	this.Data["operators"] = operatorMaps

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
	Description      string

	Must *actions.Must
}) {
	params.Must.
		Field("prefix", params.Prefix).
		Require("请选择参数")


	if len(params.Value) > 4096 {
		this.FailField("value", "对比值内容长度不能超过4096个字符")
		return
	}

	var rule = &firewallconfigs.HTTPFirewallRule{
		Id:   params.RuleId,
		IsOn: true,
	}
	if len(params.Param) > 0 {
		rule.Param = "${" + params.Prefix + "." + params.Param + "}"
	} else {
		rule.Param = "${" + params.Prefix + "}"
	}

	var paramFilters = []*firewallconfigs.ParamFilter{}
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
	rule.Description = params.Description

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
