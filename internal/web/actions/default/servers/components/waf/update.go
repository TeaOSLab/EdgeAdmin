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
	"github.com/iwind/TeaGo/types"
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
	firewallPolicy, err := dao.SharedHTTPFirewallPolicyDAO.FindEnabledHTTPFirewallPolicyConfig(this.AdminContext(), params.FirewallPolicyId)
	if err != nil {
		this.ErrorPage(err)
		return
	}
	if firewallPolicy == nil {
		this.NotFound("firewallPolicy", params.FirewallPolicyId)
		return
	}

	// block options
	if firewallPolicy.BlockOptions == nil {
		firewallPolicy.BlockOptions = firewallconfigs.NewHTTPFirewallBlockAction()
	}

	// page options
	if firewallPolicy.PageOptions == nil {
		firewallPolicy.PageOptions = firewallconfigs.NewHTTPFirewallPageAction()
	}

	// jscookie options
	if firewallPolicy.JSCookieOptions == nil {
		firewallPolicy.JSCookieOptions = firewallconfigs.NewHTTPFirewallJavascriptCookieAction()
	}

	// mode
	if len(firewallPolicy.Mode) == 0 {
		firewallPolicy.Mode = firewallconfigs.FirewallModeDefend
	}
	this.Data["modes"] = firewallconfigs.FindAllFirewallModes()

	// syn flood
	if firewallPolicy.SYNFlood == nil {
		firewallPolicy.SYNFlood = &firewallconfigs.SYNFloodConfig{
			IsOn:           false,
			MinAttempts:    10,
			TimeoutSeconds: 600,
			IgnoreLocal:    true,
		}
	}

	// log
	if firewallPolicy.Log == nil {
		firewallPolicy.Log = firewallconfigs.DefaultHTTPFirewallPolicyLogConfig
	}

	this.Data["firewallPolicy"] = maps.Map{
		"id":                 firewallPolicy.Id,
		"name":               firewallPolicy.Name,
		"description":        firewallPolicy.Description,
		"isOn":               firewallPolicy.IsOn,
		"mode":               firewallPolicy.Mode,
		"blockOptions":       firewallPolicy.BlockOptions,
		"pageOptions":        firewallPolicy.PageOptions,
		"captchaOptions":     firewallPolicy.CaptchaOptions,
		"jsCookieOptions":    firewallPolicy.JSCookieOptions,
		"useLocalFirewall":   firewallPolicy.UseLocalFirewall,
		"synFloodConfig":     firewallPolicy.SYNFlood,
		"log":                firewallPolicy.Log,
		"maxRequestBodySize": types.String(firewallPolicy.MaxRequestBodySize),
		"denyCountryHTML":    firewallPolicy.DenyCountryHTML,
		"denyProvinceHTML":   firewallPolicy.DenyProvinceHTML,
	}

	// 预置分组
	var groups = []maps.Map{}
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
	FirewallPolicyId    int64
	Name                string
	GroupCodes          []string
	BlockOptionsJSON    []byte
	PageOptionsJSON     []byte
	CaptchaOptionsJSON  []byte
	JsCookieOptionsJSON []byte
	Description         string
	IsOn                bool
	Mode                string
	UseLocalFirewall    bool
	SynFloodJSON        []byte
	LogJSON             []byte
	MaxRequestBodySize  int64
	DenyCountryHTML     string
	DenyProvinceHTML    string

	Must *actions.Must
}) {
	// 日志
	defer this.CreateLogInfo(codes.WAFPolicy_LogUpdateWAFPolicy, params.FirewallPolicyId)

	params.Must.
		Field("name", params.Name).
		Require("请输入策略名称")

	// 校验拦截选项JSON
	var blockOptions = firewallconfigs.NewHTTPFirewallBlockAction()
	err := json.Unmarshal(params.BlockOptionsJSON, blockOptions)
	if err != nil {
		this.Fail("拦截动作参数校验失败：" + err.Error())
		return
	}

	// 校验显示页面选项JSON
	var pageOptions = firewallconfigs.NewHTTPFirewallPageAction()
	err = json.Unmarshal(params.PageOptionsJSON, pageOptions)
	if err != nil {
		this.Fail("校验显示页面动作配置失败：" + err.Error())
		return
	}
	if pageOptions.Status < 100 && pageOptions.Status > 999 {
		this.Fail("显示页面动作的状态码配置错误：" + types.String(pageOptions.Status))
		return
	}

	// 校验验证码选项JSON
	var captchaOptions = firewallconfigs.NewHTTPFirewallCaptchaAction()
	err = json.Unmarshal(params.CaptchaOptionsJSON, captchaOptions)
	if err != nil {
		this.Fail("验证码动作参数校验失败：" + err.Error())
		return
	}

	// 检查极验配置
	if captchaOptions.CaptchaType == firewallconfigs.CaptchaTypeGeeTest || captchaOptions.GeeTestConfig.IsOn {
		if captchaOptions.CaptchaType == firewallconfigs.CaptchaTypeGeeTest && !captchaOptions.GeeTestConfig.IsOn {
			this.Fail("人机识别动作配置的默认验证方式为极验-行为验，所以需要选择允许用户使用极验")
			return
		}

		if len(captchaOptions.GeeTestConfig.CaptchaId) == 0 {
			this.FailField("geetestCaptchaId", "请输入极验-验证ID")
			return
		}
		if len(captchaOptions.GeeTestConfig.CaptchaKey) == 0 {
			this.FailField("geetestCaptchaKey", "请输入极验-验证Key")
			return
		}
	}

	// 校验JSCookie选项JSON
	var jsCookieOptions = firewallconfigs.NewHTTPFirewallJavascriptCookieAction()
	if len(params.JsCookieOptionsJSON) > 0 {
		err = json.Unmarshal(params.JsCookieOptionsJSON, jsCookieOptions)
		if err != nil {
			this.Fail("JSCookie动作参数校验失败：" + err.Error())
			return
		}
	}

	// 最大内容尺寸
	if params.MaxRequestBodySize < 0 {
		params.MaxRequestBodySize = 0
	}

	_, err = this.RPC().HTTPFirewallPolicyRPC().UpdateHTTPFirewallPolicy(this.AdminContext(), &pb.UpdateHTTPFirewallPolicyRequest{
		HttpFirewallPolicyId: params.FirewallPolicyId,
		IsOn:                 params.IsOn,
		Name:                 params.Name,
		Description:          params.Description,
		FirewallGroupCodes:   params.GroupCodes,
		BlockOptionsJSON:     params.BlockOptionsJSON,
		PageOptionsJSON:      params.PageOptionsJSON,
		CaptchaOptionsJSON:   params.CaptchaOptionsJSON,
		JsCookieOptionsJSON:  params.JsCookieOptionsJSON,
		Mode:                 params.Mode,
		UseLocalFirewall:     params.UseLocalFirewall,
		SynFloodJSON:         params.SynFloodJSON,
		LogJSON:              params.LogJSON,
		MaxRequestBodySize:   params.MaxRequestBodySize,
		DenyCountryHTML:      params.DenyCountryHTML,
		DenyProvinceHTML:     params.DenyProvinceHTML,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
