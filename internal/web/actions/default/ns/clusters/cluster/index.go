// Copyright 2021 Liuxiangchao iwind.liu@gmail.com. All rights reserved.

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
	"time"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.Nav("", "node", "index")
}

func (this *IndexAction) RunGet(params struct {
	ClusterId      int64
	InstalledState int
	ActiveState    int
	Keyword        string
}) {
	this.Data["installState"] = params.InstalledState
	this.Data["activeState"] = params.ActiveState
	this.Data["keyword"] = params.Keyword

	countAllResp, err := this.RPC().NSNodeRPC().CountAllEnabledNSNodesMatch(this.AdminContext(), &pb.CountAllEnabledNSNodesMatchRequest{
		NsClusterId: params.ClusterId,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	this.Data["countAll"] = countAllResp.Count

	countResp, err := this.RPC().NSNodeRPC().CountAllEnabledNSNodesMatch(this.AdminContext(), &pb.CountAllEnabledNSNodesMatchRequest{
		NsClusterId:  params.ClusterId,
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

	nodesResp, err := this.RPC().NSNodeRPC().ListEnabledNSNodesMatch(this.AdminContext(), &pb.ListEnabledNSNodesMatchRequest{
		Offset:       page.Offset,
		Size:         page.Size,
		NsClusterId:  params.ClusterId,
		InstallState: types.Int32(params.InstalledState),
		ActiveState:  types.Int32(params.ActiveState),
		Keyword:      params.Keyword,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	nodeMaps := []maps.Map{}
	for _, node := range nodesResp.NsNodes {
		// 状态
		status := &nodeconfigs.NodeStatus{}
		if len(node.StatusJSON) > 0 {
			err = json.Unmarshal(node.StatusJSON, &status)
			if err != nil {
				logs.Error(err)
				continue
			}
			status.IsActive = status.IsActive && time.Now().Unix()-status.UpdatedAt <= 60 // N秒之内认为活跃
		}

		// IP
		ipAddressesResp, err := this.RPC().NodeIPAddressRPC().FindAllEnabledNodeIPAddressesWithNodeId(this.AdminContext(), &pb.FindAllEnabledNodeIPAddressesWithNodeIdRequest{
			NodeId: node.Id,
			Role:   nodeconfigs.NodeRoleDNS,
		})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		ipAddresses := []maps.Map{}
		for _, addr := range ipAddressesResp.NodeIPAddresses {
			ipAddresses = append(ipAddresses, maps.Map{
				"id":        addr.Id,
				"name":      addr.Name,
				"ip":        addr.Ip,
				"canAccess": addr.CanAccess,
				"isOn":      addr.IsOn,
				"isUp":      addr.IsUp,
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
				"isActive":     node.IsActive,
				"updatedAt":    status.UpdatedAt,
				"hostname":     status.Hostname,
				"cpuUsage":     status.CPUUsage,
				"cpuUsageText": fmt.Sprintf("%.2f%%", status.CPUUsage*100),
				"memUsage":     status.MemoryUsage,
				"memUsageText": fmt.Sprintf("%.2f%%", status.MemoryUsage*100),
			},
			"ipAddresses": ipAddresses,
		})
	}
	this.Data["nodes"] = nodeMaps

	// 记录最近访问
	_, err = this.RPC().LatestItemRPC().IncreaseLatestItem(this.AdminContext(), &pb.IncreaseLatestItemRequest{
		ItemType: "nsCluster",
		ItemId:   params.ClusterId,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Show()
}
