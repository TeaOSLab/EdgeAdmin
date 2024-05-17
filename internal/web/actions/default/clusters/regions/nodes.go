// Copyright 2022 GoEdge CDN goedge.cdn@gmail.com. All rights reserved. Official site: https://goedge.cn .

package regions

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/maps"
)

type NodesAction struct {
	actionutils.ParentAction
}

func (this *NodesAction) Init() {
	this.Nav("", "", "node")
}

func (this *NodesAction) RunGet(params struct {
	RegionId int64
}) {
	this.Data["regionId"] = params.RegionId

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

	// 节点数量
	countResp, err := this.RPC().NodeRPC().CountAllNodeRegionInfo(this.AdminContext(), &pb.CountAllNodeRegionInfoRequest{NodeRegionId: params.RegionId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	var page = this.NewPage(countResp.Count)
	this.Data["page"] = page.AsHTML()

	// 节点列表
	var hasNodesWithoutRegion = false
	nodesResp, err := this.RPC().NodeRPC().ListNodeRegionInfo(this.AdminContext(), &pb.ListNodeRegionInfoRequest{
		NodeRegionId: params.RegionId,
		Offset:       page.Offset,
		Size:         page.Size,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	var nodeMaps = []maps.Map{}
	for _, node := range nodesResp.InfoList {
		// region
		var regionMap maps.Map
		if node.NodeRegion != nil {
			regionMap = maps.Map{
				"id":   node.NodeRegion.Id,
				"name": node.NodeRegion.Name,
			}
		} else {
			hasNodesWithoutRegion = true
		}

		// cluster
		var clusterMap maps.Map
		if node.NodeCluster != nil {
			clusterMap = maps.Map{
				"id":   node.NodeCluster.Id,
				"name": node.NodeCluster.Name,
			}
		}

		nodeMaps = append(nodeMaps, maps.Map{
			"id":      node.Id,
			"name":    node.Name,
			"region":  regionMap,
			"cluster": clusterMap,
		})
	}
	this.Data["nodes"] = nodeMaps
	this.Data["hasNodesWithoutRegion"] = hasNodesWithoutRegion

	this.Show()
}
