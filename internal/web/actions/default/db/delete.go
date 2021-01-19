package db

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
	// 创建日志
	defer this.CreateLog(oplogs.LevelInfo, "删除数据库节点 %d", params.NodeId)

	_, err := this.RPC().DBNodeRPC().DeleteDBNode(this.AdminContext(), &pb.DeleteDBNodeRequest{DbNodeId: params.NodeId})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
