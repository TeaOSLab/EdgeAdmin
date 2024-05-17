// Copyright 2021 GoEdge CDN goedge.cdn@gmail.com. All rights reserved.

package metrics

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/servers/metrics/metricutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/maps"
)

type ItemAction struct {
	actionutils.ParentAction
}

func (this *ItemAction) Init() {
	this.Nav("", "", "item")
}

func (this *ItemAction) RunGet(params struct {
	ItemId int64
}) {
	_, err := metricutils.InitItem(this.Parent(), params.ItemId)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	// 使用此指标的集群
	clustersResp, err := this.RPC().NodeClusterMetricItemRPC().FindAllNodeClustersWithMetricItemId(this.AdminContext(), &pb.FindAllNodeClustersWithMetricItemIdRequest{MetricItemId: params.ItemId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	var clusterMaps = []maps.Map{}
	for _, cluster := range clustersResp.NodeClusters {
		clusterMaps = append(clusterMaps, maps.Map{
			"id":   cluster.Id,
			"name": cluster.Name,
		})
	}
	this.Data["clusters"] = clusterMaps

	this.Show()
}
