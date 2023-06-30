package nodes

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
)

type DeleteAction struct {
	actionutils.ParentAction
}

func (this *DeleteAction) RunPost(params struct {
	ClusterId int64
	NodeId    int64
}) {
	// 创建日志
	defer this.CreateLogInfo(codes.Node_LogDeleteNodeFromCluster, params.ClusterId, params.NodeId)

	_, err := this.RPC().NodeRPC().DeleteNodeFromNodeCluster(this.AdminContext(), &pb.DeleteNodeFromNodeClusterRequest{
		NodeId:        params.NodeId,
		NodeClusterId: params.ClusterId,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
