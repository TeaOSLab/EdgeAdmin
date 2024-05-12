package providers

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/maps"
	timeutil "github.com/iwind/TeaGo/utils/time"
)

type ProviderAction struct {
	actionutils.ParentAction
}

func (this *ProviderAction) Init() {
	this.Nav("", "", "")
}

func (this *ProviderAction) RunGet(params struct {
	ProviderId int64
	Page       int
	Filter     string
}) {
	this.Data["pageNo"] = params.Page
	this.Data["filter"] = params.Filter

	providerResp, err := this.RPC().DNSProviderRPC().FindEnabledDNSProvider(this.AdminContext(), &pb.FindEnabledDNSProviderRequest{
		DnsProviderId: params.ProviderId,
		MaskParams:    true,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	var provider = providerResp.DnsProvider
	if provider == nil {
		this.NotFound("dnsProvider", params.ProviderId)
		return
	}

	var apiParams = maps.Map{}
	if len(provider.ApiParamsJSON) > 0 {
		err = json.Unmarshal(provider.ApiParamsJSON, &apiParams)
		if err != nil {
			this.ErrorPage(err)
			return
		}
	}

	// 本地EdgeDNS相关
	localEdgeDNSMap, err := this.readEdgeDNS(provider, apiParams)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Data["provider"] = maps.Map{
		"id":           provider.Id,
		"name":         provider.Name,
		"type":         provider.Type,
		"typeName":     provider.TypeName,
		"minTTL":       provider.MinTTL,
		"apiParams":    apiParams,
		"localEdgeDNS": localEdgeDNSMap,
	}

	// 域名数量
	countDomainsResp, err := this.RPC().DNSDomainRPC().CountAllDNSDomainsWithDNSProviderId(this.AdminContext(), &pb.CountAllDNSDomainsWithDNSProviderIdRequest{
		DnsProviderId: params.ProviderId,
		IsDeleted:     params.Filter == "deleted",
		IsDown:        params.Filter == "down",
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	var countDomains = countDomainsResp.Count
	var page = this.NewPage(countDomains)
	this.Data["page"] = page.AsHTML()

	// 域名
	domainsResp, err := this.RPC().DNSDomainRPC().ListBasicDNSDomainsWithDNSProviderId(this.AdminContext(), &pb.ListBasicDNSDomainsWithDNSProviderIdRequest{
		DnsProviderId: params.ProviderId,
		IsDeleted:     params.Filter == "deleted",
		IsDown:        params.Filter == "down",
		Offset:        page.Offset,
		Size:          page.Size,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	var domainMaps = []maps.Map{}
	for _, domain := range domainsResp.DnsDomains {
		dataUpdatedTime := ""
		if domain.DataUpdatedAt > 0 {
			dataUpdatedTime = timeutil.FormatTime("Y-m-d H:i:s", domain.DataUpdatedAt)
		}
		domainMaps = append(domainMaps, maps.Map{
			"id":                 domain.Id,
			"name":               domain.Name,
			"isOn":               domain.IsOn,
			"isUp":               domain.IsUp,
			"isDeleted":          domain.IsDeleted,
			"dataUpdatedTime":    dataUpdatedTime,
			"countRoutes":        len(domain.Routes),
			"countServerRecords": domain.CountServerRecords,
			"serversChanged":     domain.ServersChanged,
			"countNodeRecords":   domain.CountNodeRecords,
			"nodesChanged":       domain.NodesChanged,
			"countClusters":      domain.CountNodeClusters,
			"countAllNodes":      domain.CountAllNodes,
			"countAllServers":    domain.CountAllServers,
		})
	}
	this.Data["domains"] = domainMaps

	this.Show()
}
