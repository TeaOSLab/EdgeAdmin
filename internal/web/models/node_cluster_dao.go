package models

import (
	"context"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
)

var SharedNodeClusterDAO = new(NodeClusterDAO)

type NodeClusterDAO struct {
	BaseDAO
}

// 查找集群
func (this *NodeClusterDAO) FindEnabledNodeCluster(ctx context.Context, clusterId int64) (*pb.NodeCluster, error) {
	clusterResp, err := this.RPC().NodeClusterRPC().FindEnabledNodeCluster(ctx, &pb.FindEnabledNodeClusterRequest{NodeClusterId: clusterId})
	if err != nil {
		return nil, err
	}
	return clusterResp.NodeCluster, nil
}
