package nodes

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/oplogs"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
)

type DeleteAction struct {
	actionutils.ParentAction
}

func (this *DeleteAction) RunPost(params struct {
	NodeId int64
}) {
	// TODO 检查权限

	_, err := this.RPC().AuthorityNodeRPC().DeleteAuthorityNode(this.AdminContext(), &pb.DeleteAuthorityNodeRequest{AuthorityNodeId: params.NodeId})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	// 创建日志
	defer this.CreateLog(oplogs.LevelInfo, "删除认证节点 %d", params.NodeId)

	this.Success()
}
