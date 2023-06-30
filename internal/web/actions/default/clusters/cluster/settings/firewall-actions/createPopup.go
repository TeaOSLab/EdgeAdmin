package firewallActions

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs/firewallconfigs"
	"github.com/iwind/TeaGo/actions"
)

type CreatePopupAction struct {
	actionutils.ParentAction
}

func (this *CreatePopupAction) Init() {
	this.Nav("", "", "")
}

func (this *CreatePopupAction) RunGet(params struct {
	ClusterId int64
}) {
	this.Data["clusterId"] = params.ClusterId
	this.Data["actionTypes"] = firewallconfigs.FindAllFirewallActionTypes()

	this.Show()
}

func (this *CreatePopupAction) RunPost(params struct {
	ClusterId  int64
	Name       string
	EventLevel string
	Type       string

	// ipset
	IpsetWhiteName          string
	IpsetBlackName          string
	IpsetWhiteNameIPv6      string
	IpsetBlackNameIPv6      string
	IpsetAutoAddToIPTables  bool
	IpsetAutoAddToFirewalld bool

	// script
	ScriptPath string

	// http api
	HttpAPIURL string

	// html
	HtmlContent string

	Must *actions.Must
	CSRF *actionutils.CSRF
}) {
	defer this.CreateLogInfo(codes.WAFAction_LogCreateWAFAction, params.ClusterId)

	params.Must.
		Field("name", params.Name).
		Require("请输入动作名称").
		Field("type", params.Type).
		Require("请选择动作类型")

	var actionParams interface{} = nil
	switch params.Type {
	case firewallconfigs.FirewallActionTypeIPSet:
		params.Must.
			Field("ipsetWhiteName", params.IpsetWhiteName).
			Require("请输入IPSet白名单名称").
			Match(`^\w+$`, "请输入正确的IPSet白名单名称").
			Field("ipsetBlackName", params.IpsetBlackName).
			Require("请输入IPSet黑名单名称").
			Match(`^\w+$`, "请输入正确的IPSet黑名单名称").
			Field("ipsetWhiteNameIPv6", params.IpsetWhiteNameIPv6).
			Require("请输入IPSet IPv6白名单名称").
			Match(`^\w+$`, "请输入正确的IPSet IPv6白名单名称").
			Field("ipsetBlackNameIPv6", params.IpsetBlackNameIPv6).
			Require("请输入IPSet IPv6黑名单名称").
			Match(`^\w+$`, "请输入正确的IPSet IPv6黑名单名称")

		actionParams = &firewallconfigs.FirewallActionIPSetConfig{
			WhiteName:          params.IpsetWhiteName,
			BlackName:          params.IpsetBlackName,
			WhiteNameIPv6:      params.IpsetWhiteNameIPv6,
			BlackNameIPv6:      params.IpsetBlackNameIPv6,
			AutoAddToIPTables:  params.IpsetAutoAddToIPTables,
			AutoAddToFirewalld: params.IpsetAutoAddToFirewalld,
		}
	case firewallconfigs.FirewallActionTypeIPTables:
		actionParams = &firewallconfigs.FirewallActionIPTablesConfig{}
	case firewallconfigs.FirewallActionTypeFirewalld:
		actionParams = &firewallconfigs.FirewallActionFirewalldConfig{}
	case firewallconfigs.FirewallActionTypeScript:
		params.Must.
			Field("scriptPath", params.ScriptPath).
			Require("请输入脚本路径")
		actionParams = &firewallconfigs.FirewallActionScriptConfig{
			Path: params.ScriptPath,
		}
	case firewallconfigs.FirewallActionTypeHTTPAPI:
		params.Must.
			Field("httpAPIURL", params.HttpAPIURL).
			Require("请输入API URL").
			Match(`^(http|https):`, "API地址必须以http://或https://开头")
		actionParams = &firewallconfigs.FirewallActionHTTPAPIConfig{
			URL: params.HttpAPIURL,
		}
	case firewallconfigs.FirewallActionTypeHTML:
		params.Must.
			Field("htmlContent", params.HtmlContent).
			Require("请输入HTML内容")
		actionParams = &firewallconfigs.FirewallActionHTMLConfig{
			Content: params.HtmlContent,
		}
	default:
		this.Fail("选择的类型'" + params.Type + "'暂时不支持")
	}

	actionParamsJSON, err := json.Marshal(actionParams)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	_, err = this.RPC().NodeClusterFirewallActionRPC().CreateNodeClusterFirewallAction(this.AdminContext(), &pb.CreateNodeClusterFirewallActionRequest{
		NodeClusterId: params.ClusterId,
		Name:          params.Name,
		EventLevel:    params.EventLevel,
		Type:          params.Type,
		ParamsJSON:    actionParamsJSON,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	this.Success()
}
