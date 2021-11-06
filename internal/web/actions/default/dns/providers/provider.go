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

	// 本地EdgeDNS相关
	var localEdgeDNSMap = maps.Map{}
	if provider.Type == "localEdgeDNS" {
		nsClusterId := apiParams.GetInt64("clusterId")
		nsClusterResp, err := this.RPC().NSClusterRPC().FindEnabledNSCluster(this.AdminContext(), &pb.FindEnabledNSClusterRequest{NsClusterId: nsClusterId})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		nsCluster := nsClusterResp.NsCluster
		if nsCluster != nil {
			localEdgeDNSMap = maps.Map{
				"id":   nsCluster.Id,
				"name": nsCluster.Name,
			}
		}
	}

	this.Data["provider"] = maps.Map{
		"id":           provider.Id,
		"name":         provider.Name,
		"type":         provider.Type,
		"typeName":     provider.TypeName,
		"apiParams":    apiParams,
		"localEdgeDNS": localEdgeDNSMap,
	}

	// 域名
	domainsResp, err := this.RPC().DNSDomainRPC().FindAllEnabledDNSDomainsWithDNSProviderId(this.AdminContext(), &pb.FindAllEnabledDNSDomainsWithDNSProviderIdRequest{DnsProviderId: provider.Id})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	domainMaps := []maps.Map{}
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
