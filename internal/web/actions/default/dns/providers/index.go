package providers

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/maps"
	timeutil "github.com/iwind/TeaGo/utils/time"
	"regexp"
	"strings"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.Nav("", "", "")
}

func (this *IndexAction) RunGet(params struct {
	Keyword      string
	Domain       string
	ProviderType string
}) {
	this.Data["keyword"] = params.Keyword
	this.Data["domain"] = params.Domain
	this.Data["providerType"] = params.ProviderType

	// 格式化域名
	var domain = params.Domain
	domain = regexp.MustCompile(`^(www\.)`).ReplaceAllString(domain, "")
	domain = strings.ToLower(domain)

	countResp, err := this.RPC().DNSProviderRPC().CountAllEnabledDNSProviders(this.AdminContext(), &pb.CountAllEnabledDNSProvidersRequest{
		AdminId: this.AdminId(),
		Keyword: params.Keyword,
		Domain:  domain,
		Type:    params.ProviderType,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	count := countResp.Count
	page := this.NewPage(count)
	this.Data["page"] = page.AsHTML()

	providersResp, err := this.RPC().DNSProviderRPC().ListEnabledDNSProviders(this.AdminContext(), &pb.ListEnabledDNSProvidersRequest{
		AdminId: this.AdminId(),
		Keyword: params.Keyword,
		Domain:  domain,
		Type:    params.ProviderType,
		Offset:  page.Offset,
		Size:    page.Size,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	var providerMaps = []maps.Map{}
	for _, provider := range providersResp.DnsProviders {
		var dataUpdatedTime = ""
		if provider.DataUpdatedAt > 0 {
			dataUpdatedTime = timeutil.FormatTime("Y-m-d H:i:s", provider.DataUpdatedAt)
		}

		// 域名
		countDomainsResp, err := this.RPC().DNSDomainRPC().CountAllDNSDomainsWithDNSProviderId(this.AdminContext(), &pb.CountAllDNSDomainsWithDNSProviderIdRequest{
			DnsProviderId: provider.Id,
		})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		var countDomains = countDomainsResp.Count

		providerMaps = append(providerMaps, maps.Map{
			"id":              provider.Id,
			"name":            provider.Name,
			"type":            provider.Type,
			"typeName":        provider.TypeName,
			"dataUpdatedTime": dataUpdatedTime,
			"countDomains":    countDomains,
		})
	}
	this.Data["providers"] = providerMaps

	// 类型
	typesResponse, err := this.RPC().DNSProviderRPC().FindAllDNSProviderTypes(this.AdminContext(), &pb.FindAllDNSProviderTypesRequest{})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	var providerTypeMaps = []maps.Map{}
	for _, providerType := range typesResponse.ProviderTypes {
		countProvidersWithTypeResp, err := this.RPC().DNSProviderRPC().CountAllEnabledDNSProviders(this.AdminContext(), &pb.CountAllEnabledDNSProvidersRequest{
			Type: providerType.Code,
		})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		if countProvidersWithTypeResp.Count > 0 {
			providerTypeMaps = append(providerTypeMaps, maps.Map{
				"name":  providerType.Name,
				"code":  providerType.Code,
				"count": countProvidersWithTypeResp.Count,
			})
		}
	}
	this.Data["providerTypes"] = providerTypeMaps

	this.Show()
}
