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
			"name": t.Name,
			"code": t.Code,
		})
	}
	this.Data["types"] = typeMaps

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

	// DNS.COM
	ParamApiKey    string
	ParamApiSecret string

	Must *actions.Must
	CSRF *actionutils.CSRF
}) {
	this.CreateLog(oplogs.LevelInfo, "修改DNS服务商 %d", params.ProviderId)

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
	case "dnscom":
		params.Must.
			Field("paramApiKey", params.ParamApiKey).
			Require("请输入ApiKey").
			Field("paramApiSecret", params.ParamApiSecret).
			Require("请输入ApiSecret")

		apiParams["apiKey"] = params.ParamApiKey
		apiParams["apiSecret"] = params.ParamApiSecret
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
