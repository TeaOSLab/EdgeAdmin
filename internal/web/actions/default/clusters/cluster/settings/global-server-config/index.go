// Copyright 2022 GoEdge CDN goedge.cdn@gmail.com. All rights reserved. Official site: https://goedge.cn .

package globalServerConfig

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/types"
	"regexp"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.Nav("", "setting", "")
	this.SecondMenu("globalServerConfig")
}

func (this *IndexAction) RunGet(params struct {
	ClusterId int64
}) {
	configResp, err := this.RPC().NodeClusterRPC().FindNodeClusterGlobalServerConfig(this.AdminContext(), &pb.FindNodeClusterGlobalServerConfigRequest{NodeClusterId: params.ClusterId})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	var configJSON = configResp.GlobalServerConfigJSON
	var config = serverconfigs.NewGlobalServerConfig()
	if len(configJSON) > 0 {
		err = json.Unmarshal(configJSON, config)
		if err != nil {
			this.ErrorPage(err)
			return
		}
	}
	this.Data["config"] = config

	var httpAllDomainMismatchActionCode = serverconfigs.DomainMismatchActionPage
	var httpAllDomainMismatchActionContentHTML string
	var httpAllDomainMismatchActionStatusCode = "404"

	var httpAllDomainMismatchActionRedirectURL = ""

	if config.HTTPAll.DomainMismatchAction != nil {
		httpAllDomainMismatchActionCode = config.HTTPAll.DomainMismatchAction.Code

		if config.HTTPAll.DomainMismatchAction.Options != nil {
			// 即使是非 page 处理动作，也读取这些内容，以便于在切换到 page 时，可以顺利读取到先前的设置
			httpAllDomainMismatchActionContentHTML = config.HTTPAll.DomainMismatchAction.Options.GetString("contentHTML")
			var statusCode = config.HTTPAll.DomainMismatchAction.Options.GetInt("statusCode")
			if statusCode > 0 {
				httpAllDomainMismatchActionStatusCode = types.String(statusCode)
			}

			if config.HTTPAll.DomainMismatchAction.Code == serverconfigs.DomainMismatchActionRedirect {
				httpAllDomainMismatchActionRedirectURL = config.HTTPAll.DomainMismatchAction.Options.GetString("url")
			}
		}
	} else {
		httpAllDomainMismatchActionContentHTML = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8"/>
<title>404 not found</title>
<style>
	* { font-family: Roboto, system-ui, sans-serif; }
h3, p { text-align: center; }
p { color: grey; }
</style>
</head>
<body>
<h3>Error: 404 Page Not Found</h3>
<h3>找不到您要访问的页面。</h3>

<p>原因：找不到当前访问域名对应的网站，请联系网站管理员。</p>

</body>
</html>`
	}

	this.Data["httpAllDomainMismatchActionCode"] = httpAllDomainMismatchActionCode
	this.Data["httpAllDomainMismatchActionContentHTML"] = httpAllDomainMismatchActionContentHTML
	this.Data["httpAllDomainMismatchActionStatusCode"] = httpAllDomainMismatchActionStatusCode

	this.Data["httpAllDomainMismatchActionRedirectURL"] = httpAllDomainMismatchActionRedirectURL

	this.Show()
}

func (this *IndexAction) RunPost(params struct {
	ClusterId int64

	HttpAllMatchDomainStrictly             bool
	HttpAllDomainMismatchActionCode        string
	HttpAllDomainMismatchActionContentHTML string
	HttpAllDomainMismatchActionStatusCode  string
	HttpAllDomainMismatchActionRedirectURL string
	HttpAllAllowMismatchDomainsJSON        []byte
	HttpAllAllowNodeIP                     bool
	HttpAllDefaultDomain                   string
	HttpAllNodeIPPageHTML                  string
	HttpAllNodeIPShowPage                  bool
	HttpAllEnableServerAddrVariable        bool
	HttpAllRequestOriginsWithEncodings     bool
	HttpAllXFFMaxAddresses                 int

	HttpAllDomainAuditingIsOn   bool
	HttpAllDomainAuditingPrompt string

	HttpAllServerName                string
	HttpAllSupportsLowVersionHTTP    bool
	HttpAllMatchCertFromAllServers   bool
	HttpAllForceLnRequest            bool
	HttpAllLnRequestSchedulingMethod string

	HttpAccessLogIsOn                     bool
	HttpAccessLogEnableRequestHeaders     bool
	HttpAccessLogEnableResponseHeaders    bool
	HttpAccessLogCommonRequestHeadersOnly bool
	HttpAccessLogEnableCookies            bool
	HttpAccessLogEnableServerNotFound     bool

	LogRecordServerError bool

	PerformanceAutoReadTimeout  bool
	PerformanceAutoWriteTimeout bool
	PerformanceDebug            bool

	// TCP端口设置
	TcpAllPortRangeMin int
	TcpAllPortRangeMax int
	TcpAllDenyPorts    []int

	Must *actions.Must
	CSRF *actionutils.CSRF
}) {
	defer this.CreateLogInfo(codes.ServerGlobalSetting_LogUpdateClusterGlobalServerConfig, params.ClusterId)

	configResp, err := this.RPC().NodeClusterRPC().FindNodeClusterGlobalServerConfig(this.AdminContext(), &pb.FindNodeClusterGlobalServerConfigRequest{NodeClusterId: params.ClusterId})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	var configJSON = configResp.GlobalServerConfigJSON
	var config = serverconfigs.NewGlobalServerConfig()
	if len(configJSON) > 0 {
		err = json.Unmarshal(configJSON, config)
		if err != nil {
			this.ErrorPage(err)
			return
		}
	}

	var domainMisMatchStatusCodeString = params.HttpAllDomainMismatchActionStatusCode
	if !regexp.MustCompile(`^\d{3}$`).MatchString(domainMisMatchStatusCodeString) {
		this.FailField("httpAllDomainMismatchActionContentStatusCode", "请输入正确的状态码")
		return
	}
	var domainMisMatchStatusCode = types.Int(domainMisMatchStatusCodeString)

	config.HTTPAll.MatchDomainStrictly = params.HttpAllMatchDomainStrictly

	// validate
	if config.HTTPAll.MatchDomainStrictly {
		// validate redirect
		if params.HttpAllDomainMismatchActionCode == serverconfigs.DomainMismatchActionRedirect {
			if len(params.HttpAllDomainMismatchActionRedirectURL) == 0 {
				this.FailField("httpAllDomainMismatchActionRedirectURL", "请输入跳转目标网址URL")
				return
			}
			if !regexp.MustCompile(`(?i)(http|https)://`).MatchString(params.HttpAllDomainMismatchActionRedirectURL) {
				this.FailField("httpAllDomainMismatchActionRedirectURL", "目标网址URL必须以http://或https://开头")
				return
			}
		}
	}

	config.HTTPAll.DomainMismatchAction = &serverconfigs.DomainMismatchAction{
		Code: params.HttpAllDomainMismatchActionCode,
		Options: maps.Map{
			"statusCode":  domainMisMatchStatusCode,                      // page
			"contentHTML": params.HttpAllDomainMismatchActionContentHTML, // page
			"url":         params.HttpAllDomainMismatchActionRedirectURL, // redirect
		},
	}

	var allowMismatchDomains = []string{}
	if len(params.HttpAllAllowMismatchDomainsJSON) > 0 {
		err = json.Unmarshal(params.HttpAllAllowMismatchDomainsJSON, &allowMismatchDomains)
		if err != nil {
			this.ErrorPage(err)
			return
		}
	}

	// 域名
	config.HTTPAll.AllowMismatchDomains = allowMismatchDomains
	config.HTTPAll.AllowNodeIP = params.HttpAllAllowNodeIP
	config.HTTPAll.DefaultDomain = params.HttpAllDefaultDomain
	config.HTTPAll.NodeIPShowPage = params.HttpAllNodeIPShowPage
	config.HTTPAll.NodeIPPageHTML = params.HttpAllNodeIPPageHTML

	config.HTTPAll.DomainAuditingIsOn = params.HttpAllDomainAuditingIsOn
	config.HTTPAll.DomainAuditingPrompt = params.HttpAllDomainAuditingPrompt

	// HTTP All
	config.HTTPAll.ServerName = params.HttpAllServerName
	config.HTTPAll.SupportsLowVersionHTTP = params.HttpAllSupportsLowVersionHTTP
	config.HTTPAll.MatchCertFromAllServers = params.HttpAllMatchCertFromAllServers
	config.HTTPAll.ForceLnRequest = params.HttpAllForceLnRequest
	config.HTTPAll.LnRequestSchedulingMethod = params.HttpAllLnRequestSchedulingMethod
	config.HTTPAll.EnableServerAddrVariable = params.HttpAllEnableServerAddrVariable
	config.HTTPAll.RequestOriginsWithEncodings = params.HttpAllRequestOriginsWithEncodings
	config.HTTPAll.XFFMaxAddresses = params.HttpAllXFFMaxAddresses

	// 访问日志
	config.HTTPAccessLog.IsOn = params.HttpAccessLogIsOn
	config.HTTPAccessLog.EnableRequestHeaders = params.HttpAccessLogEnableRequestHeaders
	config.HTTPAccessLog.EnableResponseHeaders = params.HttpAccessLogEnableResponseHeaders
	config.HTTPAccessLog.CommonRequestHeadersOnly = params.HttpAccessLogCommonRequestHeadersOnly
	config.HTTPAccessLog.EnableCookies = params.HttpAccessLogEnableCookies
	config.HTTPAccessLog.EnableServerNotFound = params.HttpAccessLogEnableServerNotFound

	// 日志
	config.Log.RecordServerError = params.LogRecordServerError

	// TCP
	if params.TcpAllPortRangeMin < 1024 {
		params.TcpAllPortRangeMin = 1024
	}
	if params.TcpAllPortRangeMax > 65534 {
		params.TcpAllPortRangeMax = 65534
	} else if params.TcpAllPortRangeMax < 1024 {
		params.TcpAllPortRangeMax = 1024
	}
	if params.TcpAllPortRangeMin > params.TcpAllPortRangeMax {
		params.TcpAllPortRangeMin, params.TcpAllPortRangeMax = params.TcpAllPortRangeMax, params.TcpAllPortRangeMin
	}

	config.TCPAll.DenyPorts = params.TcpAllDenyPorts
	config.TCPAll.PortRangeMin = params.TcpAllPortRangeMin
	config.TCPAll.PortRangeMax = params.TcpAllPortRangeMax

	// 性能
	config.Performance.AutoReadTimeout = params.PerformanceAutoReadTimeout
	config.Performance.AutoWriteTimeout = params.PerformanceAutoWriteTimeout
	config.Performance.Debug = params.PerformanceDebug

	err = config.Init()
	if err != nil {
		this.Fail("配置校验失败：" + err.Error())
		return
	}

	configJSON, err = json.Marshal(config)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	_, err = this.RPC().NodeClusterRPC().UpdateNodeClusterGlobalServerConfig(this.AdminContext(), &pb.UpdateNodeClusterGlobalServerConfigRequest{
		NodeClusterId:          params.ClusterId,
		GlobalServerConfigJSON: configJSON,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
