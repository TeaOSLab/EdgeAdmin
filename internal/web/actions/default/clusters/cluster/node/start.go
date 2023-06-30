package node

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
)

type StartAction struct {
	actionutils.ParentAction
}

func (this *StartAction) RunPost(params struct {
	NodeId int64
}) {
	resp, err := this.RPC().NodeRPC().StartNode(this.AdminContext(), &pb.StartNodeRequest{NodeId: params.NodeId})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	// 创建日志
	defer this.CreateLogInfo(codes.Node_LogStartNodeRemotely, params.NodeId)

	if resp.IsOk {
		this.Success()
	}

	this.Fail("启动失败：" + resp.Error)
}
