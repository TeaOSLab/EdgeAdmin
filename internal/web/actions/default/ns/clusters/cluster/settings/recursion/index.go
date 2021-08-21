// Copyright 2021 Liuxiangchao iwind.liu@gmail.com. All rights reserved.

package recursion

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/dnsconfigs"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/actions"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.Nav("", "setting", "")
	this.SecondMenu("recursion")
}

func (this *IndexAction) RunGet(params struct {
	ClusterId int64
}) {
	this.Data["clusterId"] = params.ClusterId

	resp, err := this.RPC().NSClusterRPC().FindNSClusterRecursionConfig(this.AdminContext(), &pb.FindNSClusterRecursionConfigRequest{NsClusterId: params.ClusterId})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	var config = &dnsconfigs.RecursionConfig{}
	if len(resp.RecursionJSON) > 0 {
		err = json.Unmarshal(resp.RecursionJSON, config)
		if err != nil {
			this.ErrorPage(err)
			return
		}
	} else {
		config.UseLocalHosts = true
	}
	this.Data["config"] = config

	this.Show()
}

func (this *IndexAction) RunPost(params struct {
	ClusterId     int64
	RecursionJSON []byte

	Must *actions.Must
	CSRF *actionutils.CSRF
}) {
	defer this.CreateLogInfo("修改DNS集群 %d 的递归DNS设置", params.ClusterId)

	// TODO 校验域名

	_, err := this.RPC().NSClusterRPC().UpdateNSClusterRecursionConfig(this.AdminContext(), &pb.UpdateNSClusterRecursionConfigRequest{
		NsClusterId:   params.ClusterId,
		RecursionJSON: params.RecursionJSON,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	this.Success()
}
