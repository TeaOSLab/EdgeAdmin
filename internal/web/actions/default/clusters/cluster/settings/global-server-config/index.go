// Copyright 2022 Liuxiangchao iwind.liu@gmail.com. All rights reserved. Official site: https://goedge.cn .

package globalServerConfig

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
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
	var config = serverconfigs.DefaultGlobalServerConfig()
	if len(configJSON) > 0 {
		err = json.Unmarshal(configJSON, config)
		if err != nil {
			this.ErrorPage(err)
			return
		}
	}
	this.Data["config"] = config

	var httpAllDomainMismatchActionContentHTML = ""
	if config.HTTPAll.DomainMismatchAction != nil {
		httpAllDomainMismatchActionContentHTML = config.HTTPAll.DomainMismatchAction.Options.GetString("contentHTML")
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
	this.Data["httpAllDomainMismatchActionContentHTML"] = httpAllDomainMismatchActionContentHTML

	this.Show()
}

func (this *IndexAction) RunPost(params struct {
	ClusterId int64

	HttpAllMatchDomainStrictly             bool
	HttpAllDomainMismatchActionContentHTML string
	HttpAllAllowMismatchDomainsJSON        []byte
	HttpAllAllowNodeIP                     bool
	HttpAllDefaultDomain                   string

	HttpAccessLogEnableRequestHeaders     bool
	HttpAccessLogEnableResponseHeaders    bool
	HttpAccessLogCommonRequestHeadersOnly bool
	HttpAccessLogEnableCookies            bool

	LogRecordServerError bool

	Must *actions.Must
	CSRF *actionutils.CSRF
}) {
	defer this.CreateLogInfo("修改集群 %d 全局配置", params.ClusterId)

	configResp, err := this.RPC().NodeClusterRPC().FindNodeClusterGlobalServerConfig(this.AdminContext(), &pb.FindNodeClusterGlobalServerConfigRequest{NodeClusterId: params.ClusterId})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	var configJSON = configResp.GlobalServerConfigJSON
	var config = serverconfigs.DefaultGlobalServerConfig()
	if len(configJSON) > 0 {
		err = json.Unmarshal(configJSON, config)
		if err != nil {
			this.ErrorPage(err)
			return
		}
	}

	config.HTTPAll.MatchDomainStrictly = params.HttpAllMatchDomainStrictly
	config.HTTPAll.DomainMismatchAction = &serverconfigs.DomainMismatchAction{
		Code: serverconfigs.DomainMismatchActionPage,
		Options: maps.Map{
			"statusCode":  404,
			"contentHTML": params.HttpAllDomainMismatchActionContentHTML,
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

	config.HTTPAll.AllowMismatchDomains = allowMismatchDomains
	config.HTTPAll.AllowNodeIP = params.HttpAllAllowNodeIP
	config.HTTPAll.DefaultDomain = params.HttpAllDefaultDomain

	config.HTTPAccessLog.EnableRequestHeaders = params.HttpAccessLogEnableRequestHeaders
	config.HTTPAccessLog.EnableResponseHeaders = params.HttpAccessLogEnableResponseHeaders
	config.HTTPAccessLog.CommonRequestHeadersOnly = params.HttpAccessLogCommonRequestHeadersOnly
	config.HTTPAccessLog.EnableCookies = params.HttpAccessLogEnableCookies

	config.Log.RecordServerError = params.LogRecordServerError

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
