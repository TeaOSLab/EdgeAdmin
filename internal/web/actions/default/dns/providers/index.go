package providers

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/maps"
	timeutil "github.com/iwind/TeaGo/utils/time"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.Nav("", "", "")
}

func (this *IndexAction) RunGet(params struct {
	Keyword string
}) {
	this.Data["keyword"] = params.Keyword

	countResp, err := this.RPC().DNSProviderRPC().CountAllEnabledDNSProviders(this.AdminContext(), &pb.CountAllEnabledDNSProvidersRequest{
		AdminId: this.AdminId(),
		Keyword: params.Keyword,
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
		Offset:  page.Offset,
		Size:    page.Size,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	providerMaps := []maps.Map{}
	for _, provider := range providersResp.DnsProviders {
		dataUpdatedTime := ""
		if provider.DataUpdatedAt > 0 {
			dataUpdatedTime = timeutil.FormatTime("Y-m-d H:i:s", provider.DataUpdatedAt)
		}

		// 域名
		countDomainsResp, err := this.RPC().DNSDomainRPC().CountAllEnabledDNSDomainsWithDNSProviderId(this.AdminContext(), &pb.CountAllEnabledDNSDomainsWithDNSProviderIdRequest{DnsProviderId: provider.Id})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		countDomains := countDomainsResp.Count

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

	this.Show()
}
