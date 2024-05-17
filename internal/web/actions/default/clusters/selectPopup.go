// Copyright 2021 GoEdge CDN goedge.cdn@gmail.com. All rights reserved.

package clusters

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/utils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/maps"
)

type SelectPopupAction struct {
	actionutils.ParentAction
}

func (this *SelectPopupAction) Init() {
	this.Nav("", "", "")
}

func (this *SelectPopupAction) RunGet(params struct {
	SelectedClusterIds string
	Keyword            string
	PageSize           int64
}) {
	this.Data["keyword"] = params.Keyword

	var selectedIds = utils.SplitNumbers(params.SelectedClusterIds)

	countResp, err := this.RPC().NodeClusterRPC().CountAllEnabledNodeClusters(this.AdminContext(), &pb.CountAllEnabledNodeClustersRequest{
		Keyword: params.Keyword,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	var count = countResp.Count
	var pageSize int64 = 6
	if params.PageSize > 0 {
		pageSize = params.PageSize
	}
	var page = this.NewPage(count, pageSize)
	this.Data["page"] = page.AsHTML()

	clustersResp, err := this.RPC().NodeClusterRPC().ListEnabledNodeClusters(this.AdminContext(), &pb.ListEnabledNodeClustersRequest{
		Keyword: params.Keyword,
		Offset:  page.Offset,
		Size:    page.Size,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	var clusterMaps = []maps.Map{}
	for _, cluster := range clustersResp.NodeClusters {
		// 节点数
		countNodesResp, err := this.RPC().NodeRPC().CountAllEnabledNodesMatch(this.AdminContext(), &pb.CountAllEnabledNodesMatchRequest{NodeClusterId: cluster.Id})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		var countNodes = countNodesResp.Count

		clusterMaps = append(clusterMaps, maps.Map{
			"id":         cluster.Id,
			"name":       cluster.Name,
			"isOn":       cluster.IsOn,
			"countNodes": countNodes,
			"isSelected": lists.ContainsInt64(selectedIds, cluster.Id),
		})
	}
	this.Data["clusters"] = clusterMaps

	this.Show()
}
