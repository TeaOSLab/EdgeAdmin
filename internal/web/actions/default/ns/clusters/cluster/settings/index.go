// Copyright 2021 Liuxiangchao iwind.liu@gmail.com. All rights reserved.

package cluster

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.Nav("", "setting", "")
	this.SecondMenu("basic")
}

func (this *IndexAction) RunGet(params struct {
	ClusterId int64
}) {
	clusterResp, err := this.RPC().NSClusterRPC().FindEnabledNSCluster(this.AdminContext(), &pb.FindEnabledNSClusterRequest{NsClusterId: params.ClusterId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	cluster := clusterResp.NsCluster
	if cluster == nil {
		this.NotFound("nsCluster", params.ClusterId)
		return
	}

	this.Data["cluster"] = maps.Map{
		"id":   cluster.Id,
		"name": cluster.Name,
		"isOn": cluster.IsOn,
	}

	this.Show()
}

func (this *IndexAction) RunPost(params struct {
	ClusterId int64
	Name      string
	IsOn      bool

	Must *actions.Must
	CSRF *actionutils.CSRF
}) {
	defer this.CreateLogInfo("修改域名服务集群基本信息 %d", params.ClusterId)

	params.Must.
		Field("name", params.Name).
		Require("请输入集群名称")

	_, err := this.RPC().NSClusterRPC().UpdateNSCluster(this.AdminContext(), &pb.UpdateNSClusterRequest{
		NsClusterId: params.ClusterId,
		Name:        params.Name,
		IsOn:        params.IsOn,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
