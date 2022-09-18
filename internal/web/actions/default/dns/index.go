package dns

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/maps"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.Nav("dns", "dns", "")
}

func (this *IndexAction) RunGet(params struct {
	Keyword string
}) {
	this.Data["keyword"] = params.Keyword

	countResp, err := this.RPC().NodeClusterRPC().CountAllEnabledNodeClusters(this.AdminContext(), &pb.CountAllEnabledNodeClustersRequest{
		Keyword: params.Keyword,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	page := this.NewPage(countResp.Count)
	this.Data["page"] = page.AsHTML()

	clustersResp, err := this.RPC().NodeClusterRPC().ListEnabledNodeClusters(this.AdminContext(), &pb.ListEnabledNodeClustersRequest{
		Keyword: params.Keyword,
		Offset:  page.Offset,
		Size:    page.Size,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	clusterMaps := []maps.Map{}
	for _, cluster := range clustersResp.NodeClusters {
		domainId := cluster.DnsDomainId
		domainName := ""
		providerId := int64(0)
		providerName := ""
		providerTypeName := ""

		if cluster.DnsDomainId > 0 {
			domainResp, err := this.RPC().DNSDomainRPC().FindBasicDNSDomain(this.AdminContext(), &pb.FindBasicDNSDomainRequest{DnsDomainId: domainId})
			if err != nil {
				this.ErrorPage(err)
				return
			}
			domain := domainResp.DnsDomain
			if domain == nil {
				domainId = 0
			} else {
				domainName = domain.Name
				providerResp, err := this.RPC().DNSProviderRPC().FindEnabledDNSProvider(this.AdminContext(), &pb.FindEnabledDNSProviderRequest{DnsProviderId: domain.ProviderId})
				if err != nil {
					this.ErrorPage(err)
					return
				}
				if providerResp.DnsProvider != nil {
					providerId = providerResp.DnsProvider.Id
					providerName = providerResp.DnsProvider.Name
					providerTypeName = providerResp.DnsProvider.TypeName
				}
			}
		}

		clusterMaps = append(clusterMaps, maps.Map{
			"id":               cluster.Id,
			"name":             cluster.Name,
			"dnsName":          cluster.DnsName,
			"domainId":         domainId,
			"domainName":       domainName,
			"providerId":       providerId,
			"providerName":     providerName,
			"providerTypeName": providerTypeName,
		})
	}
	this.Data["clusters"] = clusterMaps

	this.Show()
}
