// Copyright 2023 GoEdge CDN goedge.cdn@gmail.com. All rights reserved. Official site: https://goedge.cn .

package clusters

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/types"
)

type CreateNodeAction struct {
	actionutils.ParentAction
}

func (this *CreateNodeAction) Init() {
	this.Nav("", "cluster", "createNode")
}

func (this *CreateNodeAction) RunGet(params struct{}) {
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

	// 如果只有一个默认集群，那么直接跳转到集群
	clustersResp, err := this.RPC().NodeClusterRPC().ListEnabledNodeClusters(this.AdminContext(), &pb.ListEnabledNodeClustersRequest{
		Offset:  0,
		Size:    2,
		Keyword: "",
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	if len(clustersResp.NodeClusters) == 1 {
		this.RedirectURL("/clusters/cluster/createNode?clusterId=" + types.String(clustersResp.NodeClusters[0].Id))
		return
	}

	this.Show()
}
