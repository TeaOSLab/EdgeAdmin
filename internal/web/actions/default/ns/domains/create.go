// Copyright 2021 Liuxiangchao iwind.liu@gmail.com. All rights reserved.

package domains

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/dns/domains/domainutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/actions"
	"strings"
)

type CreateAction struct {
	actionutils.ParentAction
}

func (this *CreateAction) Init() {
	this.Nav("", "", "create")
}

func (this *CreateAction) RunGet(params struct{}) {
	// 集群数量
	countClustersResp, err := this.RPC().NSClusterRPC().CountAllEnabledNSClusters(this.AdminContext(), &pb.CountAllEnabledNSClustersRequest{})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	this.Data["countClusters"] = countClustersResp.Count

	this.Show()
}

func (this *CreateAction) RunPost(params struct {
	Name      string
	ClusterId int64
	UserId    int64

	Must *actions.Must
	CSRF *actionutils.CSRF
}) {
	var domainId int64
	defer func() {
		this.CreateLogInfo("创建域名 %d", domainId)
	}()

	params.Name = strings.ToLower(params.Name)

	params.Must.
		Field("name", params.Name).
		Require("请输入域名").
		Expect(func() (message string, success bool) {
			success = domainutils.ValidateDomainFormat(params.Name)
			if !success {
				message = "请输入正确的域名"
			}
			return
		}).
		Field("clusterId", params.ClusterId).
		Gt(0, "请选择所属集群")

	createResp, err := this.RPC().NSDomainRPC().CreateNSDomain(this.AdminContext(), &pb.CreateNSDomainRequest{
		NsClusterId: params.ClusterId,
		UserId:      params.UserId,
		Name:        params.Name,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	domainId = createResp.NsDomainId

	this.Success()
}
