package clusters

import (
	"encoding/json"
	"fmt"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/configutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/nodeconfigs"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/types"
	"time"
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
	isSearching := len(params.Keyword) > 0
	this.Data["keyword"] = params.Keyword
	this.Data["searchType"] = params.SearchType
	this.Data["isSearching"] = isSearching

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

	// 搜索节点
	if params.SearchType == "node" && len(params.Keyword) > 0 {
		this.searchNodes(params.Keyword)
		return
	}

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

func (this *IndexAction) searchNodes(keyword string) {
	// 搜索节点
	countResp, err := this.RPC().NodeRPC().CountAllEnabledNodesMatch(this.AdminContext(), &pb.CountAllEnabledNodesMatchRequest{
		Keyword: keyword,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	count := countResp.Count
	page := this.NewPage(count)
	this.Data["page"] = page.AsHTML()
	this.Data["countNodes"] = count

	nodesResp, err := this.RPC().NodeRPC().ListEnabledNodesMatch(this.AdminContext(), &pb.ListEnabledNodesMatchRequest{
		Offset:  page.Offset,
		Size:    page.Size,
		Keyword: keyword,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	nodeMaps := []maps.Map{}
	for _, node := range nodesResp.Nodes {
		// 状态
		isSynced := false
		status := &nodeconfigs.NodeStatus{}
		if len(node.StatusJSON) > 0 {
			err = json.Unmarshal(node.StatusJSON, &status)
			if err != nil {
				this.ErrorPage(err)
				return
			}
			status.IsActive = status.IsActive && time.Now().Unix()-status.UpdatedAt <= 60 // N秒之内认为活跃
			isSynced = status.ConfigVersion == node.Version
		}

		// IP
		ipAddressesResp, err := this.RPC().NodeIPAddressRPC().FindAllEnabledIPAddressesWithNodeId(this.AdminContext(), &pb.FindAllEnabledIPAddressesWithNodeIdRequest{
			NodeId: node.Id,
			Role:   nodeconfigs.NodeRoleNode,
		})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		ipAddresses := []maps.Map{}
		for _, addr := range ipAddressesResp.Addresses {
			ipAddresses = append(ipAddresses, maps.Map{
				"id":        addr.Id,
				"name":      addr.Name,
				"ip":        addr.Ip,
				"canAccess": addr.CanAccess,
			})
		}

		// 分组
		var groupMap maps.Map = nil
		if node.NodeGroup != nil {
			groupMap = maps.Map{
				"id":   node.NodeGroup.Id,
				"name": node.NodeGroup.Name,
			}
		}

		// 区域
		var regionMap maps.Map = nil
		if node.NodeRegion != nil {
			regionMap = maps.Map{
				"id":   node.NodeRegion.Id,
				"name": node.NodeRegion.Name,
			}
		}

		// DNS
		dnsRouteNames := []string{}
		for _, route := range node.DnsRoutes {
			dnsRouteNames = append(dnsRouteNames, route.Name)
		}

		// 从集群
		var secondaryClusterMaps []maps.Map
		for _, secondaryCluster := range node.SecondaryNodeClusters {
			secondaryClusterMaps = append(secondaryClusterMaps, maps.Map{
				"id":   secondaryCluster.Id,
				"name": secondaryCluster.Name,
				"isOn": secondaryCluster.IsOn,
			})
		}

		nodeMaps = append(nodeMaps, maps.Map{
			"id":          node.Id,
			"name":        node.Name,
			"isInstalled": node.IsInstalled,
			"isOn":        node.IsOn,
			"isUp":        node.IsUp,
			"installStatus": maps.Map{
				"isRunning":  node.InstallStatus.IsRunning,
				"isFinished": node.InstallStatus.IsFinished,
				"isOk":       node.InstallStatus.IsOk,
				"error":      node.InstallStatus.Error,
			},
			"status": maps.Map{
				"isActive":     status.IsActive,
				"updatedAt":    status.UpdatedAt,
				"hostname":     status.Hostname,
				"cpuUsage":     status.CPUUsage,
				"cpuUsageText": fmt.Sprintf("%.2f%%", status.CPUUsage*100),
				"memUsage":     status.MemoryUsage,
				"memUsageText": fmt.Sprintf("%.2f%%", status.MemoryUsage*100),
			},
			"cluster": maps.Map{
				"id":   node.NodeCluster.Id,
				"name": node.NodeCluster.Name,
			},
			"secondaryClusters": secondaryClusterMaps,
			"isSynced":          isSynced,
			"ipAddresses":       ipAddresses,
			"group":             groupMap,
			"region":            regionMap,
			"dnsRouteNames":     dnsRouteNames,
		})
	}
	this.Data["nodes"] = nodeMaps

	this.Data["clusters"] = []maps.Map{}

	// 搜索集群
	{
		countResp, err := this.RPC().NodeClusterRPC().CountAllEnabledNodeClusters(this.AdminContext(), &pb.CountAllEnabledNodeClustersRequest{
			Keyword: keyword,
		})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		this.Data["countClusters"] = countResp.Count
	}

	this.Show()
}
