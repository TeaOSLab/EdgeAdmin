// Copyright 2021 Liuxiangchao iwind.liu@gmail.com. All rights reserved.

package clusters

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/actions"
)

type CreateAction struct {
	actionutils.ParentAction
}

func (this *CreateAction) Init() {
	this.Nav("", "", "create")
}

func (this *CreateAction) RunGet(params struct{}) {
	this.Show()
}

func (this *CreateAction) RunPost(params struct {
	Name string

	Must *actions.Must
	CSRF *actionutils.CSRF
}) {
	var clusterId int64
	defer this.CreateLogInfo("创建域名服务集群 %d", clusterId)

	params.Must.
		Field("name", params.Name).
		Require("请输入集群名称")

	resp, err := this.RPC().NSClusterRPC().CreateNSCluster(this.AdminContext(), &pb.CreateNSClusterRequest{
		Name: params.Name,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	clusterId = resp.NsClusterId

	this.Success()
}
