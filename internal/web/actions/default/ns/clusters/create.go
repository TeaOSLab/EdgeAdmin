// Copyright 2021 Liuxiangchao iwind.liu@gmail.com. All rights reserved.

package clusters

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/dnsconfigs"
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
	// 默认的访问日志设置
	this.Data["accessLogRef"] = &dnsconfigs.NSAccessLogRef{
		IsOn: true,
	}

	this.Show()
}

func (this *CreateAction) RunPost(params struct {
	Name          string
	AccessLogJSON []byte

	Must *actions.Must
	CSRF *actionutils.CSRF
}) {
	var clusterId int64
	defer func() {
		this.CreateLogInfo("创建域名服务集群 %d", clusterId)
	}()

	params.Must.
		Field("name", params.Name).
		Require("请输入集群名称")

	// 校验访问日志设置
	ref := &dnsconfigs.NSAccessLogRef{}
	err := json.Unmarshal(params.AccessLogJSON, ref)
	if err != nil {
		this.Fail("数据格式错误：" + err.Error())
	}
	err = ref.Init()
	if err != nil {
		this.Fail("数据格式错误：" + err.Error())
	}

	resp, err := this.RPC().NSClusterRPC().CreateNSCluster(this.AdminContext(), &pb.CreateNSClusterRequest{
		Name:          params.Name,
		AccessLogJSON: params.AccessLogJSON,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	clusterId = resp.NsClusterId

	this.Success()
}
