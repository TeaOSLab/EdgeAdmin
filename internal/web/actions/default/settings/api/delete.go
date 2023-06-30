package api

import (	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
)

type DeleteAction struct {
	actionutils.ParentAction
}

func (this *DeleteAction) RunPost(params struct {
	NodeId int64
}) {
	// 创建日志
	defer this.CreateLogInfo(codes.APINode_LogDeleteAPINode, params.NodeId)

	// 检查是否是唯一的节点
	nodeResp, err := this.RPC().APINodeRPC().FindEnabledAPINode(this.AdminContext(), &pb.FindEnabledAPINodeRequest{ApiNodeId: params.NodeId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	var apiNode = nodeResp.ApiNode
	if apiNode == nil {
		this.Success()
		return
	}
	if apiNode.IsOn {
		countResp, err := this.RPC().APINodeRPC().CountAllEnabledAndOnAPINodes(this.AdminContext(), &pb.CountAllEnabledAndOnAPINodesRequest{})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		if countResp.Count == 1 {
			this.Fail("无法删除此节点：必须保留至少一个可用的API节点")
		}
	}

	_, err = this.RPC().APINodeRPC().DeleteAPINode(this.AdminContext(), &pb.DeleteAPINodeRequest{ApiNodeId: params.NodeId})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
