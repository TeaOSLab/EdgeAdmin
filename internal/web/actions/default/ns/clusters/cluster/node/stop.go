package node

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/oplogs"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
)

type StopAction struct {
	actionutils.ParentAction
}

func (this *StopAction) RunPost(params struct {
	NodeId int64
}) {
	resp, err := this.RPC().NSNodeRPC().StopNSNode(this.AdminContext(), &pb.StopNSNodeRequest{NsNodeId: params.NodeId})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	// 创建日志
	defer this.CreateLog(oplogs.LevelInfo, "远程停止节点 %d", params.NodeId)

	if resp.IsOk {
		this.Success()
	}

	this.Fail("执行失败：" + resp.Error)
}
