// Copyright 2022 GoEdge CDN goedge.cdn@gmail.com. All rights reserved.

package clusters

import (
	"encoding/json"
	teaconst "github.com/TeaOSLab/EdgeAdmin/internal/const"
	"github.com/TeaOSLab/EdgeAdmin/internal/utils/numberutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/nodeconfigs"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/types"
	"strconv"
	"time"
)

type NodesAction struct {
	actionutils.ParentAction
}

func (this *NodesAction) Init() {
	this.Nav("", "", "node")
}

func (this *NodesAction) RunGet(params struct {
	ClusterId      int64
	GroupId        int64
	RegionId       int64
	InstalledState int
	ActiveState    int
	Keyword        string
	Level          int32

	CpuOrder         string
	MemoryOrder      string
	TrafficInOrder   string
	TrafficOutOrder  string
	LoadOrder        string
	ConnectionsOrder string
}) {
	this.Data["groupId"] = params.GroupId
	this.Data["regionId"] = params.RegionId
	this.Data["installState"] = params.InstalledState
	this.Data["activeState"] = params.ActiveState
	this.Data["keyword"] = params.Keyword
	this.Data["level"] = params.Level
	this.Data["clusterId"] = params.ClusterId
	this.Data["hasOrder"] = len(params.CpuOrder) > 0 || len(params.MemoryOrder) > 0 || len(params.TrafficInOrder) > 0 || len(params.TrafficOutOrder) > 0 || len(params.LoadOrder) > 0 || len(params.ConnectionsOrder) > 0

	// 集群是否已经设置了线路
	clusterDNSResp, err := this.RPC().NodeClusterRPC().FindEnabledNodeClusterDNS(this.AdminContext(), &pb.FindEnabledNodeClusterDNSRequest{NodeClusterId: params.ClusterId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	this.Data["hasClusterDNS"] = clusterDNSResp.Domain != nil

	// 数量
	countAllResp, err := this.RPC().NodeRPC().CountAllEnabledNodesMatch(this.AdminContext(), &pb.CountAllEnabledNodesMatchRequest{
		NodeClusterId: params.ClusterId,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	this.Data["countAll"] = countAllResp.Count

	countResp, err := this.RPC().NodeRPC().CountAllEnabledNodesMatch(this.AdminContext(), &pb.CountAllEnabledNodesMatchRequest{
		NodeClusterId: params.ClusterId,
		NodeGroupId:   params.GroupId,
		NodeRegionId:  params.RegionId,
		Level:         params.Level,
		InstallState:  types.Int32(params.InstalledState),
		ActiveState:   types.Int32(params.ActiveState),
		Keyword:       params.Keyword,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	var page = this.NewPage(countResp.Count)
	this.Data["page"] = page.AsHTML()

	var req = &pb.ListEnabledNodesMatchRequest{
		Offset:        page.Offset,
		Size:          page.Size,
		NodeClusterId: params.ClusterId,
		NodeGroupId:   params.GroupId,
		NodeRegionId:  params.RegionId,
		Level:         params.Level,
		InstallState:  types.Int32(params.InstalledState),
		ActiveState:   types.Int32(params.ActiveState),
		Keyword:       params.Keyword,
	}
	if params.CpuOrder == "asc" {
		req.CpuAsc = true
	} else if params.CpuOrder == "desc" {
		req.CpuDesc = true
	} else if params.MemoryOrder == "asc" {
		req.MemoryAsc = true
	} else if params.MemoryOrder == "desc" {
		req.MemoryDesc = true
	} else if params.TrafficInOrder == "asc" {
		req.TrafficInAsc = true
	} else if params.TrafficInOrder == "desc" {
		req.TrafficInDesc = true
	} else if params.TrafficOutOrder == "asc" {
		req.TrafficOutAsc = true
	} else if params.TrafficOutOrder == "desc" {
		req.TrafficOutDesc = true
	} else if params.LoadOrder == "asc" {
		req.LoadAsc = true
	} else if params.LoadOrder == "desc" {
		req.LoadDesc = true
	} else if params.ConnectionsOrder == "asc" {
		req.ConnectionsAsc = true
	} else if params.ConnectionsOrder == "desc" {
		req.ConnectionsDesc = true
	}
	nodesResp, err := this.RPC().NodeRPC().ListEnabledNodesMatch(this.AdminContext(), req)
	if err != nil {
		this.ErrorPage(err)
		return
	}
	var nodeMaps = []maps.Map{}
	for _, node := range nodesResp.Nodes {
		// 状态
		var isSynced = false
		var status = &nodeconfigs.NodeStatus{}
		if len(node.StatusJSON) > 0 {
			err = json.Unmarshal(node.StatusJSON, &status)
			if err != nil {
				logs.Error(err)
				continue
			}
			status.IsActive = status.IsActive && time.Now().Unix()-status.UpdatedAt <= 60 // N秒之内认为活跃
			isSynced = status.ConfigVersion == node.Version
		}

		// IP
		ipAddressesResp, err := this.RPC().NodeIPAddressRPC().FindAllEnabledNodeIPAddressesWithNodeId(this.AdminContext(), &pb.FindAllEnabledNodeIPAddressesWithNodeIdRequest{
			NodeId: node.Id,
			Role:   nodeconfigs.NodeRoleNode,
		})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		var ipAddresses = []maps.Map{}
		for _, addr := range ipAddressesResp.NodeIPAddresses {
			// 专属集群
			var addrClusterMaps = []maps.Map{}
			for _, addrCluster := range addr.NodeClusters {
				addrClusterMaps = append(addrClusterMaps, maps.Map{
					"id":   addrCluster.Id,
					"name": addrCluster.Name,
				})
			}

			ipAddresses = append(ipAddresses, maps.Map{
				"id":        addr.Id,
				"name":      addr.Name,
				"ip":        addr.Ip,
				"canAccess": addr.CanAccess,
				"isUp":      addr.IsUp,
				"isOn":      addr.IsOn,
				"clusters":  addrClusterMaps,
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
		var dnsRouteNames = []string{}
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
			"isBackup":    node.IsBackupForCluster || node.IsBackupForGroup,
			"offlineDay":  node.OfflineDay,
			"installStatus": maps.Map{
				"isRunning":  node.InstallStatus.IsRunning,
				"isFinished": node.InstallStatus.IsFinished,
				"isOk":       node.InstallStatus.IsOk,
				"error":      node.InstallStatus.Error,
			},
			"status": maps.Map{
				"isActive":         status.IsActive,
				"updatedAt":        status.UpdatedAt,
				"hostname":         status.Hostname,
				"cpuUsage":         status.CPUUsage,
				"cpuUsageText":     numberutils.FormatFloat2(status.CPUUsage*100) + "%",
				"memUsage":         status.MemoryUsage,
				"memUsageText":     numberutils.FormatFloat2(status.MemoryUsage*100) + "%",
				"trafficInBytes":   status.TrafficInBytes,
				"trafficOutBytes":  status.TrafficOutBytes,
				"load1m":           numberutils.PadFloatZero(numberutils.FormatFloat2(status.Load1m), 2),
				"countConnections": status.ConnectionCount,
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
			"level":             node.Level,
		})
	}
	this.Data["nodes"] = nodeMaps

	// 所有分组
	var groupMaps = []maps.Map{}
	groupsResp, err := this.RPC().NodeGroupRPC().FindAllEnabledNodeGroupsWithNodeClusterId(this.AdminContext(), &pb.FindAllEnabledNodeGroupsWithNodeClusterIdRequest{
		NodeClusterId: params.ClusterId,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	for _, group := range groupsResp.NodeGroups {
		countNodesInGroupResp, err := this.RPC().NodeRPC().CountAllEnabledNodesWithNodeGroupId(this.AdminContext(), &pb.CountAllEnabledNodesWithNodeGroupIdRequest{NodeGroupId: group.Id})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		countNodes := countNodesInGroupResp.Count
		groupName := group.Name
		if countNodes > 0 {
			groupName += "(" + strconv.FormatInt(countNodes, 10) + ")"
		}
		groupMaps = append(groupMaps, maps.Map{
			"id":         group.Id,
			"name":       groupName,
			"countNodes": countNodes,
		})
	}
	this.Data["groups"] = groupMaps

	// 所有区域
	regionsResp, err := this.RPC().NodeRegionRPC().FindAllAvailableNodeRegions(this.AdminContext(), &pb.FindAllAvailableNodeRegionsRequest{})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	var regionMaps = []maps.Map{}
	for _, region := range regionsResp.NodeRegions {
		regionMaps = append(regionMaps, maps.Map{
			"id":   region.Id,
			"name": region.Name,
		})
	}
	this.Data["regions"] = regionMaps

	// 级别
	this.Data["levels"] = []maps.Map{}
	if teaconst.IsPlus {
		this.Data["levels"] = nodeconfigs.FindAllNodeLevels()
	}

	// 集群总数
	totalClustersResp, err := this.RPC().NodeClusterRPC().CountAllEnabledNodeClusters(this.AdminContext(), &pb.CountAllEnabledNodeClustersRequest{})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	this.Data["totalNodeClusters"] = totalClustersResp.Count

	// 节点总数
	this.Data["totalNodes"] = countResp.Count

	this.Show()
}
