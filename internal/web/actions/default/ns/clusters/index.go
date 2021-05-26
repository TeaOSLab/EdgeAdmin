// Copyright 2021 Liuxiangchao iwind.liu@gmail.com. All rights reserved.

package clusters

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/configutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/types"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.Nav("", "", "index")
}

func (this *IndexAction) RunGet(params struct{}) {
	countResp, err := this.RPC().NSClusterRPC().CountAllEnabledNSClusters(this.AdminContext(), &pb.CountAllEnabledNSClustersRequest{})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	count := countResp.Count
	page := this.NewPage(count)
	this.Data["page"] = page.AsHTML()

	clustersResp, err := this.RPC().NSClusterRPC().ListEnabledNSClusters(this.AdminContext(), &pb.ListEnabledNSClustersRequest{
		Offset: page.Offset,
		Size:   page.Size,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	clusterMaps := []maps.Map{}
	for _, cluster := range clustersResp.NsClusters {
		// 全部节点数量
		countNodesResp, err := this.RPC().NSNodeRPC().CountAllEnabledNSNodesMatch(this.AdminContext(), &pb.CountAllEnabledNSNodesMatchRequest{NsClusterId: cluster.Id})
		if err != nil {
			this.ErrorPage(err)
			return
		}

		// 在线节点
		countActiveNodesResp, err := this.RPC().NSNodeRPC().CountAllEnabledNSNodesMatch(this.AdminContext(), &pb.CountAllEnabledNSNodesMatchRequest{
			NsClusterId: cluster.Id,
			ActiveState: types.Int32(configutils.BoolStateYes),
		})
		if err != nil {
			this.ErrorPage(err)
			return
		}

		// 需要升级的节点
		countUpgradeNodesResp, err := this.RPC().NSNodeRPC().CountAllUpgradeNSNodesWithNSClusterId(this.AdminContext(), &pb.CountAllUpgradeNSNodesWithNSClusterIdRequest{NsClusterId: cluster.Id})
		if err != nil {
			this.ErrorPage(err)
			return
		}

		clusterMaps = append(clusterMaps, maps.Map{
			"id":                cluster.Id,
			"name":              cluster.Name,
			"isOn":              cluster.IsOn,
			"countAllNodes":     countNodesResp.Count,
			"countActiveNodes":  countActiveNodesResp.Count,
			"countUpgradeNodes": countUpgradeNodesResp.Count,
		})
	}
	this.Data["clusters"] = clusterMaps

	this.Show()
}
