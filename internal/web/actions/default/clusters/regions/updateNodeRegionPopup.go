// Copyright 2022 GoEdge CDN goedge.cdn@gmail.com. All rights reserved. Official site: https://goedge.cn .

package regions

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type UpdateNodeRegionPopupAction struct {
	actionutils.ParentAction
}

func (this *UpdateNodeRegionPopupAction) Init() {
	this.Nav("", "", "")
}

func (this *UpdateNodeRegionPopupAction) RunGet(params struct {
	NodeId   int64
	RegionId int64
}) {
	// node
	nodeResp, err := this.RPC().NodeRPC().FindEnabledNode(this.AdminContext(), &pb.FindEnabledNodeRequest{NodeId: params.NodeId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	var node = nodeResp.Node
	if node == nil {
		this.NotFound("node", params.NodeId)
		return
	}
	this.Data["node"] = maps.Map{
		"id":   node.Id,
		"name": node.Name,
	}

	// region
	this.Data["region"] = maps.Map{
		"id":   0,
		"name": "",
	}
	if params.RegionId > 0 {
		regionResp, err := this.RPC().NodeRegionRPC().FindEnabledNodeRegion(this.AdminContext(), &pb.FindEnabledNodeRegionRequest{NodeRegionId: params.RegionId})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		var region = regionResp.NodeRegion
		if region != nil {
			this.Data["region"] = maps.Map{
				"id":   region.Id,
				"name": region.Name,
			}
		}
	}

	// all regions
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

	this.Show()
}

func (this *UpdateNodeRegionPopupAction) RunPost(params struct {
	NodeId   int64
	RegionId int64

	Must *actions.Must
	CSRF *actionutils.CSRF
}) {
	defer this.CreateLogInfo(codes.NodeRegion_LogMoveNodeBetweenRegions, params.RegionId)

	_, err := this.RPC().NodeRPC().UpdateNodeRegionInfo(this.AdminContext(), &pb.UpdateNodeRegionInfoRequest{
		NodeId:       params.NodeId,
		NodeRegionId: params.RegionId,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
