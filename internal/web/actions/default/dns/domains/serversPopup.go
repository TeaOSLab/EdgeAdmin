package domains

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/maps"
)

type ServersPopupAction struct {
	actionutils.ParentAction
}

func (this *ServersPopupAction) Init() {
	this.Nav("", "", "")
}

func (this *ServersPopupAction) RunGet(params struct {
	DomainId int64
}) {
	this.Data["domainId"] = params.DomainId

	// 域名信息
	domainResp, err := this.RPC().DNSDomainRPC().FindBasicDNSDomain(this.AdminContext(), &pb.FindBasicDNSDomainRequest{
		DnsDomainId: params.DomainId,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	var domain = domainResp.DnsDomain
	if domain == nil {
		this.NotFound("dnsDomain", params.DomainId)
		return
	}

	this.Data["domain"] = domain.Name

	// 服务信息
	var clusterMaps = []maps.Map{}
	clustersResp, err := this.RPC().NodeClusterRPC().FindAllEnabledNodeClustersWithDNSDomainId(this.AdminContext(), &pb.FindAllEnabledNodeClustersWithDNSDomainIdRequest{DnsDomainId: params.DomainId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	for _, cluster := range clustersResp.NodeClusters {
		serversResp, err := this.RPC().ServerRPC().FindAllEnabledServersDNSWithNodeClusterId(this.AdminContext(), &pb.FindAllEnabledServersDNSWithNodeClusterIdRequest{NodeClusterId: cluster.Id})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		var serverMaps = []maps.Map{}
		for _, server := range serversResp.Servers {
			var isOk = false
			if len(cluster.DnsName) > 0 && len(server.DnsName) > 0 {
				checkResp, err := this.RPC().DNSDomainRPC().ExistDNSDomainRecord(this.AdminContext(), &pb.ExistDNSDomainRecordRequest{
					DnsDomainId: params.DomainId,
					Name:        server.DnsName,
					Type:        "CNAME",
					Value:       cluster.DnsName + "." + domain.Name,
				})
				if err != nil {
					this.ErrorPage(err)
					return
				}
				isOk = checkResp.IsOk
			}

			serverMaps = append(serverMaps, maps.Map{
				"id":      server.Id,
				"name":    server.Name,
				"dnsName": server.DnsName,
				"isOk":    isOk,
			})
		}
		clusterMaps = append(clusterMaps, maps.Map{
			"id":      cluster.Id,
			"name":    cluster.Name,
			"dnsName": cluster.DnsName,
			"servers": serverMaps,
		})
	}
	this.Data["clusters"] = clusterMaps

	this.Show()
}
