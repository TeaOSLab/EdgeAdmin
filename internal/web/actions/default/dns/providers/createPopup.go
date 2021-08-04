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

func (this *CreatePopupAction) RunPost(params struct {
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
