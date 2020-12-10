package cluster

import (
	"encoding/json"
	"fmt"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/nodeconfigs"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/types"
	"strconv"
	"time"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.Nav("", "node", "index")
	this.SecondMenu("nodes")
}

func (this *IndexAction) RunGet(params struct {
	ClusterId      int64
	GroupId        int64
	InstalledState int
	ActiveState    int
	Keyword        string
}) {
	this.Data["groupId"] = params.GroupId
	this.Data["installState"] = params.InstalledState
	this.Data["activeState"] = params.ActiveState
	this.Data["keyword"] = params.Keyword

	countResp, err := this.RPC().NodeRPC().CountAllEnabledNodesMatch(this.AdminContext(), &pb.CountAllEnabledNodesMatchRequest{
		ClusterId:    params.ClusterId,
		GroupId:      params.GroupId,
		InstallState: types.Int32(params.InstalledState),
		ActiveState:  types.Int32(params.ActiveState),
		Keyword:      params.Keyword,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	page := this.NewPage(countResp.Count)
	this.Data["page"] = page.AsHTML()

	nodesResp, err := this.RPC().NodeRPC().ListEnabledNodesMatch(this.AdminContext(), &pb.ListEnabledNodesMatchRequest{
		Offset:       page.Offset,
		Size:         page.Size,
		ClusterId:    params.ClusterId,
		GroupId:      params.GroupId,
		InstallState: types.Int32(params.InstalledState),
		ActiveState:  types.Int32(params.ActiveState),
		Keyword:      params.Keyword,
	})
	nodeMaps := []maps.Map{}
	for _, node := range nodesResp.Nodes {
		// 状态
		isSynced := false
		status := &nodeconfigs.NodeStatus{}
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
		ipAddressesResp, err := this.RPC().NodeIPAddressRPC().FindAllEnabledIPAddressesWithNodeId(this.AdminContext(), &pb.FindAllEnabledIPAddressesWithNodeIdRequest{NodeId: node.Id})
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

		var groupMap maps.Map = nil
		if node.Group != nil {
			groupMap = maps.Map{
				"id":   node.Group.Id,
				"name": node.Group.Name,
			}
		}

		// DNS
		dnsRouteNames := []string{}
		for _, route := range node.DnsRoutes {
			dnsRouteNames = append(dnsRouteNames, route.Name)
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
				"id":   node.Cluster.Id,
				"name": node.Cluster.Name,
			},
			"isSynced":      isSynced,
			"ipAddresses":   ipAddresses,
			"group":         groupMap,
			"dnsRouteNames": dnsRouteNames,
		})
	}
	this.Data["nodes"] = nodeMaps

	// 所有分组
	groupMaps := []maps.Map{}
	groupsResp, err := this.RPC().NodeGroupRPC().FindAllEnabledNodeGroupsWithClusterId(this.AdminContext(), &pb.FindAllEnabledNodeGroupsWithClusterIdRequest{
		ClusterId: params.ClusterId,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	for _, group := range groupsResp.Groups {
		countResp, err := this.RPC().NodeRPC().CountAllEnabledNodesWithNodeGroupId(this.AdminContext(), &pb.CountAllEnabledNodesWithNodeGroupIdRequest{NodeGroupId: group.Id})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		countNodes := countResp.Count
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

	this.Show()
}
