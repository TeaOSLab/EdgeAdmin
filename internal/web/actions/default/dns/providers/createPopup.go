// Copyright 2022 Liuxiangchao iwind.liu@gmail.com. All rights reserved.

//go:build !plus

package providers

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/oplogs"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/rands"
)

type CreatePopupAction struct {
	actionutils.ParentAction
}

func (this *CreatePopupAction) Init() {
	this.Nav("", "", "")
}

func (this *CreatePopupAction) RunGet(params struct{}) {
	// 所有厂商
	typesResp, err := this.RPC().DNSProviderRPC().FindAllDNSProviderTypes(this.AdminContext(), &pb.FindAllDNSProviderTypesRequest{})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	typeMaps := []maps.Map{}
	for _, t := range typesResp.ProviderTypes {
		typeMaps = append(typeMaps, maps.Map{
			"name":        t.Name,
			"code":        t.Code,
			"description": t.Description,
		})
	}
	this.Data["types"] = typeMaps

	// 自动生成CustomHTTP私钥
	this.Data["paramCustomHTTPSecret"] = rands.HexString(32)

	// EdgeDNS集群列表
	this.Data["nsClusters"] = []maps.Map{}

	this.Show()
}

func (this *CreatePopupAction) RunPost(params struct {
	Name string
	Type string

	// DNSPod
	ParamId     string
	ParamToken  string
	ParamRegion string

	// AliDNS
	ParamAliDNSAccessKeyId     string
	ParamAliDNSAccessKeySecret string
	ParamAliDNSRegionId        string

	// HuaweiDNS
	ParamHuaweiAccessKeyId     string
	ParamHuaweiAccessKeySecret string

	// CloudFlare
	ParamCloudFlareAPIKey string
	ParamCloudFlareEmail  string

	// CustomHTTP
	ParamCustomHTTPURL    string
	ParamCustomHTTPSecret string

	// EdgeDNS API
	ParamEdgeDNSAPIRole            string
	ParamEdgeDNSAPIHost            string
	ParamEdgeDNSAPIAccessKeyId     string
	ParamEdgeDNSAPIAccessKeySecret string

	Must *actions.Must
	CSRF *actionutils.CSRF
}) {
	params.Must.
		Field("name", params.Name).
		Require("请输入账号说明").
		Field("type", params.Type).
		Require("请选择服务商厂家")

	apiParams := maps.Map{}
	switch params.Type {
	case "dnspod":
		params.Must.
			Field("paramId", params.ParamId).
			Require("请输入密钥ID").
			Field("paramToken", params.ParamToken).
			Require("请输入密钥Token")

		apiParams["id"] = params.ParamId
		apiParams["token"] = params.ParamToken
		apiParams["region"] = params.ParamRegion
	case "alidns":
		params.Must.
			Field("paramAliDNSAccessKeyId", params.ParamAliDNSAccessKeyId).
			Require("请输入AccessKeyId").
			Field("paramAliDNSAccessKeySecret", params.ParamAliDNSAccessKeySecret).
			Require("请输入AccessKeySecret")

		apiParams["accessKeyId"] = params.ParamAliDNSAccessKeyId
		apiParams["accessKeySecret"] = params.ParamAliDNSAccessKeySecret
		apiParams["regionId"] = params.ParamAliDNSRegionId
	case "huaweiDNS":
		params.Must.
			Field("paramHuaweiAccessKeyId", params.ParamHuaweiAccessKeyId).
			Require("请输入AccessKeyId").
			Field("paramHuaweiAccessKeySecret", params.ParamHuaweiAccessKeySecret).
			Require("请输入AccessKeySecret")

		apiParams["accessKeyId"] = params.ParamHuaweiAccessKeyId
		apiParams["accessKeySecret"] = params.ParamHuaweiAccessKeySecret
	case "cloudFlare":
		params.Must.
			Field("paramCloudFlareAPIKey", params.ParamCloudFlareAPIKey).
			Require("请输入API密钥").
			Field("paramCloudFlareEmail", params.ParamCloudFlareEmail).
			Email("请输入正确格式的邮箱地址")
		apiParams["apiKey"] = params.ParamCloudFlareAPIKey
		apiParams["email"] = params.ParamCloudFlareEmail
	case "edgeDNSAPI":
		params.Must.
			Field("paramEdgeDNSAPIHost", params.ParamEdgeDNSAPIHost).
			Require("请输入API地址").
			Field("paramEdgeDNSAPIRole", params.ParamEdgeDNSAPIRole).
			Require("请选择AccessKey类型").
			Field("paramEdgeDNSAPIAccessKeyId", params.ParamEdgeDNSAPIAccessKeyId).
			Require("请输入AccessKey ID").
			Field("paramEdgeDNSAPIAccessKeySecret", params.ParamEdgeDNSAPIAccessKeySecret).
			Require("请输入AccessKey密钥")
		apiParams["host"] = params.ParamEdgeDNSAPIHost
		apiParams["role"] = params.ParamEdgeDNSAPIRole
		apiParams["accessKeyId"] = params.ParamEdgeDNSAPIAccessKeyId
		apiParams["accessKeySecret"] = params.ParamEdgeDNSAPIAccessKeySecret
	case "customHTTP":
		params.Must.
			Field("paramCustomHTTPURL", params.ParamCustomHTTPURL).
			Require("请输入HTTP URL").
			Match("^(?i)(http|https):", "URL必须以http://或者https://开头").
			Field("paramCustomHTTPSecret", params.ParamCustomHTTPSecret).
			Require("请输入私钥")
		apiParams["url"] = params.ParamCustomHTTPURL
		apiParams["secret"] = params.ParamCustomHTTPSecret
	default:
		this.Fail("暂时不支持此服务商'" + params.Type + "'")
	}

	createResp, err := this.RPC().DNSProviderRPC().CreateDNSProvider(this.AdminContext(), &pb.CreateDNSProviderRequest{
		Name:          params.Name,
		Type:          params.Type,
		ApiParamsJSON: apiParams.AsJSON(),
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	defer this.CreateLog(oplogs.LevelInfo, "创建DNS服务商 %d", createResp.DnsProviderId)

	this.Success()
}
