//go:build !plus

package providers

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type UpdatePopupAction struct {
	actionutils.ParentAction
}

func (this *UpdatePopupAction) Init() {
	this.Nav("", "", "")
}

func (this *UpdatePopupAction) RunGet(params struct {
	ProviderId int64
}) {
	providerResp, err := this.RPC().DNSProviderRPC().FindEnabledDNSProvider(this.AdminContext(), &pb.FindEnabledDNSProviderRequest{
		DnsProviderId: params.ProviderId,
		MaskParams:    true,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	provider := providerResp.DnsProvider
	if provider == nil {
		this.NotFound("dnsProvider", params.ProviderId)
		return
	}

	apiParams := maps.Map{}
	if len(provider.ApiParamsJSON) > 0 {
		err = json.Unmarshal(provider.ApiParamsJSON, &apiParams)
		if err != nil {
			this.ErrorPage(err)
			return
		}
	}

	this.Data["provider"] = maps.Map{
		"id":       provider.Id,
		"name":     provider.Name,
		"type":     provider.Type,
		"typeName": provider.TypeName,
		"minTTL":   provider.MinTTL,
		"params":   apiParams,
	}

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

	// EdgeDNS集群列表
	this.Data["nsClusters"] = []maps.Map{}

	this.Show()
}

func (this *UpdatePopupAction) RunPost(params struct {
	ProviderId int64

	Name string
	Type string

	// DNSPod
	ParamDNSPodId     string
	ParamDNSPodToken  string
	ParamDNSPodRegion string

	ParamDNSPodAPIType         string
	ParamDNSPodAccessKeyId     string
	ParamDNSPodAccessKeySecret string

	// AliDNS
	ParamAliDNSAccessKeyId     string
	ParamAliDNSAccessKeySecret string
	ParamAliDNSRegionId        string

	// HuaweiDNS
	ParamHuaweiAccessKeyId     string
	ParamHuaweiAccessKeySecret string
	ParamHuaweiEndpoint        string

	// CloudFlare
	ParamCloudFlareAPIKey string
	ParamCloudFlareEmail  string

	// CustomHTTP
	ParamCustomHTTPURL    string
	ParamCustomHTTPSecret string

	// EdgeDNS API
	ParamEdgeDNSAPIHost            string
	ParamEdgeDNSAPIRole            string
	ParamEdgeDNSAPIAccessKeyId     string
	ParamEdgeDNSAPIAccessKeySecret string

	MinTTL int32

	Must *actions.Must
	CSRF *actionutils.CSRF
}) {
	defer this.CreateLogInfo(codes.DNSProvider_LogUpdateDNSProvider, params.ProviderId)

	params.Must.
		Field("name", params.Name).
		Require("请输入账号说明").
		Field("type", params.Type).
		Require("请选择服务商厂家")

	var apiParams = maps.Map{}
	switch params.Type {
	case "dnspod":
		apiParams["apiType"] = params.ParamDNSPodAPIType
		switch params.ParamDNSPodAPIType {
		case "tencentDNS":
			params.Must.
				Field("paramDNSPodAccessKeyId", params.ParamDNSPodAccessKeyId).
				Require("请输入SecretId").
				Field("paramDNSPodAccessKeySecret", params.ParamDNSPodAccessKeySecret).
				Require("请输入SecretKey")
			apiParams["accessKeyId"] = params.ParamDNSPodAccessKeyId
			apiParams["accessKeySecret"] = params.ParamDNSPodAccessKeySecret
			apiParams["region"] = params.ParamDNSPodRegion
		default:
			params.Must.
				Field("paramId", params.ParamDNSPodId).
				Require("请输入密钥ID").
				Field("paramToken", params.ParamDNSPodToken).
				Require("请输入密钥Token")

			apiParams["id"] = params.ParamDNSPodId
			apiParams["token"] = params.ParamDNSPodToken
			apiParams["region"] = params.ParamDNSPodRegion
		}
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
		apiParams["endpoint"] = params.ParamHuaweiEndpoint
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

	_, err := this.RPC().DNSProviderRPC().UpdateDNSProvider(this.AdminContext(), &pb.UpdateDNSProviderRequest{
		DnsProviderId: params.ProviderId,
		Name:          params.Name,
		MinTTL:        params.MinTTL,
		ApiParamsJSON: apiParams.AsJSON(),
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
