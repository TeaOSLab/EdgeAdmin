package clusters

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/configutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/nodeconfigs"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/types"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.Nav("", "cluster", "index")
}

func (this *IndexAction) RunGet(params struct {
	Keyword    string
	SearchType string
}) {
	var isSearching = len(params.Keyword) > 0
	this.Data["keyword"] = params.Keyword
	this.Data["searchType"] = params.SearchType
	this.Data["isSearching"] = isSearching

	// 集群总数
	totalClustersResp, err := this.RPC().NodeClusterRPC().CountAllEnabledNodeClusters(this.AdminContext(), &pb.CountAllEnabledNodeClustersRequest{})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	this.Data["totalNodeClusters"] = totalClustersResp.Count

	// 节点总数
	totalNodesResp, err := this.RPC().NodeRPC().CountAllEnabledNodes(this.AdminContext(), &pb.CountAllEnabledNodesRequest{})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	this.Data["totalNodes"] = totalNodesResp.Count

	// 常用的集群
	latestClusterMaps := []maps.Map{}
	if !isSearching {
		clustersResp, err := this.RPC().NodeClusterRPC().FindLatestNodeClusters(this.AdminContext(), &pb.FindLatestNodeClustersRequest{Size: 6})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		for _, cluster := range clustersResp.NodeClusters {
			latestClusterMaps = append(latestClusterMaps, maps.Map{
				"id":   cluster.Id,
				"name": cluster.Name,
			})
		}
	}
	this.Data["latestClusters"] = latestClusterMaps

	// 搜索集群
	countResp, err := this.RPC().NodeClusterRPC().CountAllEnabledNodeClusters(this.AdminContext(), &pb.CountAllEnabledNodeClustersRequest{
		Keyword: params.Keyword,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	this.Data["countClusters"] = countResp.Count

	count := countResp.Count
	page := this.NewPage(count)
	this.Data["page"] = page.AsHTML()

	clusterMaps := []maps.Map{}
	if count > 0 {
		clustersResp, err := this.RPC().NodeClusterRPC().ListEnabledNodeClusters(this.AdminContext(), &pb.ListEnabledNodeClustersRequest{
			Keyword: params.Keyword,
			Offset:  page.Offset,
			Size:    page.Size,
		})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		for _, cluster := range clustersResp.NodeClusters {
			// 全部节点数量
			countNodesResp, err := this.RPC().NodeRPC().CountAllEnabledNodesMatch(this.AdminContext(), &pb.CountAllEnabledNodesMatchRequest{NodeClusterId: cluster.Id})
			if err != nil {
				this.ErrorPage(err)
				return
			}

			// 在线节点
			countActiveNodesResp, err := this.RPC().NodeRPC().CountAllEnabledNodesMatch(this.AdminContext(), &pb.CountAllEnabledNodesMatchRequest{
				NodeClusterId: cluster.Id,
				ActiveState:   types.Int32(configutils.BoolStateYes),
			})
			if err != nil {
				this.ErrorPage(err)
				return
			}

			// 需要升级的节点
			countUpgradeNodesResp, err := this.RPC().NodeRPC().CountAllUpgradeNodesWithNodeClusterId(this.AdminContext(), &pb.CountAllUpgradeNodesWithNodeClusterIdRequest{NodeClusterId: cluster.Id})
			if err != nil {
				this.ErrorPage(err)
				return
			}

			// DNS
			dnsDomainName := ""
			if cluster.DnsDomainId > 0 {
				dnsInfoResp, err := this.RPC().NodeClusterRPC().FindEnabledNodeClusterDNS(this.AdminContext(), &pb.FindEnabledNodeClusterDNSRequest{NodeClusterId: cluster.Id})
				if err != nil {
					this.ErrorPage(err)
					return
				}
				if dnsInfoResp.Domain != nil {
					dnsDomainName = dnsInfoResp.Domain.Name
				}
			}

			// 服务数
			countServersResp, err := this.RPC().ServerRPC().CountAllEnabledServersWithNodeClusterId(this.AdminContext(), &pb.CountAllEnabledServersWithNodeClusterIdRequest{NodeClusterId: cluster.Id})
			if err != nil {
				this.ErrorPage(err)
			}

			if cluster.TimeZone == nodeconfigs.DefaultTimeZoneLocation {
				cluster.TimeZone = ""
			}

			clusterMaps = append(clusterMaps, maps.Map{
				"id":                cluster.Id,
				"name":              cluster.Name,
				"installDir":        cluster.InstallDir,
				"countAllNodes":     countNodesResp.Count,
				"countActiveNodes":  countActiveNodesResp.Count,
				"countUpgradeNodes": countUpgradeNodesResp.Count,
				"dnsDomainId":       cluster.DnsDomainId,
				"dnsName":           cluster.DnsName,
				"dnsDomainName":     dnsDomainName,
				"countServers":      countServersResp.Count,
				"timeZone":          cluster.TimeZone,
				"isPinned":          cluster.IsPinned,
			})
		}
	}
	this.Data["clusters"] = clusterMaps
	this.Data["nodes"] = []maps.Map{}

	if len(params.Keyword) > 0 {
		// 搜索节点
		countResp, err := this.RPC().NodeRPC().CountAllEnabledNodesMatch(this.AdminContext(), &pb.CountAllEnabledNodesMatchRequest{
			Keyword: params.Keyword,
		})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		this.Data["countNodes"] = countResp.Count
	}

	this.Show()
}
