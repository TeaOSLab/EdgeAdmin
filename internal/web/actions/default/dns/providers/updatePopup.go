package providers

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/oplogs"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
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
	providerResp, err := this.RPC().DNSProviderRPC().FindEnabledDNSProvider(this.AdminContext(), &pb.FindEnabledDNSProviderRequest{DnsProviderId: params.ProviderId})
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
	nsClustersResp, err := this.RPC().NSClusterRPC().FindAllEnabledNSClusters(this.AdminContext(), &pb.FindAllEnabledNSClustersRequest{})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	nsClusterMaps := []maps.Map{}
	for _, nsCluster := range nsClustersResp.NsClusters {
		nsClusterMaps = append(nsClusterMaps, maps.Map{
			"id":   nsCluster.Id,
			"name": nsCluster.Name,
		})
	}
	this.Data["nsClusters"] = nsClusterMaps

	this.Show()
}

func (this *UpdatePopupAction) RunPost(params struct {
	ProviderId int64

	Name string
	Type string

	// DNSPod
	ParamId    string
	ParamToken string

	// AliDNS
	ParamAccessKeyId     string
	ParamAccessKeySecret string

	// HuaweiDNS
	ParamHuaweiAccessKeyId     string
	ParamHuaweiAccessKeySecret string

	// DNS.COM
	ParamApiKey    string
	ParamApiSecret string

	// CloudFlare
	CloudFlareAPIKey string
	CloudFlareEmail  string

	// Local EdgeDNS
	ParamLocalEdgeDNSClusterId int64

	// CustomHTTP
	ParamCustomHTTPURL    string
	ParamCustomHTTPSecret string

	Must *actions.Must
	CSRF *actionutils.CSRF
}) {
	defer this.CreateLog(oplogs.LevelInfo, "修改DNS服务商 %d", params.ProviderId)

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
	case "alidns":
		params.Must.
			Field("paramAccessKeyId", params.ParamAccessKeyId).
			Require("请输入AccessKeyId").
			Field("paramAccessKeySecret", params.ParamAccessKeySecret).
			Require("请输入AccessKeySecret")

		apiParams["accessKeyId"] = params.ParamAccessKeyId
		apiParams["accessKeySecret"] = params.ParamAccessKeySecret
	case "huaweiDNS":
		params.Must.
			Field("paramHuaweiAccessKeyId", params.ParamHuaweiAccessKeyId).
			Require("请输入AccessKeyId").
			Field("paramHuaweiAccessKeySecret", params.ParamHuaweiAccessKeySecret).
			Require("请输入AccessKeySecret")

		apiParams["accessKeyId"] = params.ParamHuaweiAccessKeyId
		apiParams["accessKeySecret"] = params.ParamHuaweiAccessKeySecret
	case "dnscom":
		params.Must.
			Field("paramApiKey", params.ParamApiKey).
			Require("请输入ApiKey").
			Field("paramApiSecret", params.ParamApiSecret).
			Require("请输入ApiSecret")

		apiParams["apiKey"] = params.ParamApiKey
		apiParams["apiSecret"] = params.ParamApiSecret
	case "cloudFlare":
		params.Must.
			Field("cloudFlareAPIKey", params.CloudFlareAPIKey).
			Require("请输入API密钥").
			Field("cloudFlareEmail", params.CloudFlareEmail).
			Email("请输入正确格式的邮箱地址")
		apiParams["apiKey"] = params.CloudFlareAPIKey
		apiParams["email"] = params.CloudFlareEmail
	case "localEdgeDNS":
		params.Must.
			Field("ParamLocalEdgeDNSClusterId", params.ParamLocalEdgeDNSClusterId).
			Gt(0, "请选择域名服务集群")
		apiParams["clusterId"] = params.ParamLocalEdgeDNSClusterId
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
		ApiParamsJSON: apiParams.AsJSON(),
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
