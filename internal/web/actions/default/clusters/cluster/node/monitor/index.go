// Copyright 2021 Liuxiangchao iwind.liu@gmail.com. All rights reserved.

package monitor

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/clusters/cluster/node/nodeutils"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.Nav("", "node", "monitor")
}

func (this *IndexAction) RunGet(params struct {
	NodeId int64
}) {
	this.Data["nodeId"] = params.NodeId

	// 初始化节点信息（用于菜单）
	err := nodeutils.InitNodeInfo(this, params.NodeId)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Show()
}
