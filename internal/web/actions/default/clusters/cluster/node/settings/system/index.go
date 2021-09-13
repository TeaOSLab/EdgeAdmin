// Copyright 2021 Liuxiangchao iwind.liu@gmail.com. All rights reserved.

package system

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/clusters/cluster/node/nodeutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.Nav("", "", "update")
	this.SecondMenu("system")
}

func (this *IndexAction) RunGet(params struct {
	NodeId int64
}) {
	node, err := nodeutils.InitNodeInfo(this.Parent(), params.NodeId)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	// 获取节点信息
	var nodeMap = this.Data["node"].(maps.Map)
	nodeMap["maxCPU"] = node.MaxCPU

	this.Show()
}

func (this *IndexAction) RunPost(params struct {
	NodeId int64
	MaxCPU int32

	Must *actions.Must
	CSRF *actionutils.CSRF
}) {
	defer this.CreateLogInfo("修改节点 %d 系统信息", params.NodeId)

	if params.MaxCPU < 0 {
		this.Fail("CPU线程数不能小于0")
	}

	_, err := this.RPC().NodeRPC().UpdateNodeSystem(this.AdminContext(), &pb.UpdateNodeSystemRequest{
		NodeId: params.NodeId,
		MaxCPU: params.MaxCPU,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
