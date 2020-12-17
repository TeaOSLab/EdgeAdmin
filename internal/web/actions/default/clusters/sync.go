package clusters

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/nodes/nodeutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/messageconfigs"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
)

// 同步集群
type SyncAction struct {
	actionutils.ParentAction
}

func (this *SyncAction) RunPost(params struct{}) {
	// TODO 将来可以单独选择某一个集群进行单独的同步

	// 所有有变化的集群
	clustersResp, err := this.RPC().NodeClusterRPC().FindAllChangedNodeClusters(this.AdminContext(), &pb.FindAllChangedNodeClustersRequest{})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	clusters := clustersResp.NodeClusters

	for _, cluster := range clusters {
		_, err := this.RPC().NodeRPC().SyncNodesVersionWithCluster(this.AdminContext(), &pb.SyncNodesVersionWithClusterRequest{
			NodeClusterId: cluster.Id,
		})
		if err != nil {
			this.ErrorPage(err)
			return
		}

		// 发送通知
		_, err = nodeutils.SendMessageToCluster(this.AdminContext(), cluster.Id, messageconfigs.MessageCodeConfigChanged, &messageconfigs.ConfigChangedMessage{}, 10)
		if err != nil {
			this.ErrorPage(err)
			return
		}
	}

	this.Success()
}
