package domains

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/utils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/maps"
)

type NodesPopupAction struct {
	actionutils.ParentAction
}

func (this *NodesPopupAction) Init() {
	this.Nav("", "", "")
}

func (this *NodesPopupAction) RunGet(params struct {
	DomainId int64
}) {
	// 域名信息
	domainResp, err := this.RPC().DNSDomainRPC().FindEnabledBasicDNSDomain(this.AdminContext(), &pb.FindEnabledBasicDNSDomainRequest{
		DnsDomainId: params.DomainId,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	domain := domainResp.DnsDomain
	if domain == nil {
		this.NotFound("dnsDomain", params.DomainId)
		return
	}

	this.Data["domain"] = domain.Name

	// 集群
	clusterMaps := []maps.Map{}
	clustersResp, err := this.RPC().NodeClusterRPC().FindAllEnabledNodeClustersWithDNSDomainId(this.AdminContext(), &pb.FindAllEnabledNodeClustersWithDNSDomainIdRequest{DnsDomainId: params.DomainId})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	for _, cluster := range clustersResp.NodeClusters {
		// 节点DNS解析记录
		nodesResp, err := this.RPC().NodeRPC().FindAllEnabledNodesDNSWithNodeClusterId(this.AdminContext(), &pb.FindAllEnabledNodesDNSWithNodeClusterIdRequest{NodeClusterId: cluster.Id})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		nodeMaps := []maps.Map{}
		for _, node := range nodesResp.Nodes {
			if len(node.Routes) > 0 {
				for _, route := range node.Routes {
					// 检查是否有域名解析记录
					isOk := false
					if len(route.Name) > 0 && len(node.IpAddr) > 0 && len(cluster.DnsName) > 0 {
						var recordType = "A"
						if utils.IsIPv6(node.IpAddr) {
							recordType = "AAAA"
						}
						checkResp, err := this.RPC().DNSDomainRPC().ExistDNSDomainRecord(this.AdminContext(), &pb.ExistDNSDomainRecordRequest{
							DnsDomainId: params.DomainId,
							Name:        cluster.DnsName,
							Type:        recordType,
							Route:       route.Code,
							Value:       node.IpAddr,
						})
						if err != nil {
							this.ErrorPage(err)
							return
						}
						isOk = checkResp.IsOk
					}

					nodeMaps = append(nodeMaps, maps.Map{
						"id":     node.Id,
						"name":   node.Name,
						"ipAddr": node.IpAddr,
						"route": maps.Map{
							"name": route.Name,
							"code": route.Code,
						},
						"clusterId": node.NodeClusterId,
						"isOk":      isOk,
					})
				}
			} else {
				nodeMaps = append(nodeMaps, maps.Map{
					"id":     node.Id,
					"name":   node.Name,
					"ipAddr": node.IpAddr,
					"route": maps.Map{
						"name": "",
						"code": "",
					},
					"clusterId": node.NodeClusterId,
					"isOk":      false,
				})
			}
		}

		if len(nodeMaps) == 0 {
			continue
		}

		clusterMaps = append(clusterMaps, maps.Map{
			"id":      cluster.Id,
			"name":    cluster.Name,
			"dnsName": cluster.DnsName,
			"nodes":   nodeMaps,
		})
	}
	this.Data["clusters"] = clusterMaps

	this.Show()
}
