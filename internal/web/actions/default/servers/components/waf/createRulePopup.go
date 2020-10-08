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
	for _, def := range firewallconfigs.AllCheckpoints {
		if (params.Type == "inbound" && def.IsRequest) || (params.Type == "outbound" && !def.IsRequest) {
			checkpointList = append(checkpointList, maps.Map{
				"name":        def.Name,
				"prefix":      def.Prefix,
				"description": def.Description,
				"hasParams":   len(def.Params) > 0,
				"params":      def.Params,
				"options":     def.Options,
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
	RuleId      int64
	Prefix      string
	Operator    string
	Param       string
	OptionsJSON []byte
	Value       string
	Case        bool

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
			rule.CheckpointOptions[option.GetString("code")] = option.GetString("value")
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
