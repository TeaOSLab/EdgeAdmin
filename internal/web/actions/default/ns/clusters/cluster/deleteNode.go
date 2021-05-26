// Copyright 2021 Liuxiangchao iwind.liu@gmail.com. All rights reserved.

package cluster

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
)

type DeleteNodeAction struct {
	actionutils.ParentAction
}

func (this *DeleteNodeAction) RunPost(params struct {
	NodeId int64
}) {
	defer this.CreateLogInfo("删除域名服务节点 %d", params.NodeId)

	_, err := this.RPC().NSNodeRPC().DeleteNSNode(this.AdminContext(), &pb.DeleteNSNodeRequest{NsNodeId: params.NodeId})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
