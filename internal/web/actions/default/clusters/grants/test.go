// Copyright 2021 GoEdge CDN goedge.cdn@gmail.com. All rights reserved.

package grants

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/clusters/grants/grantutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
	"strings"
)

type TestAction struct {
	actionutils.ParentAction
}

func (this *TestAction) Init() {
	this.Nav("", "", "test")
}

func (this *TestAction) RunGet(params struct {
	GrantId int64
}) {
	grantResp, err := this.RPC().NodeGrantRPC().FindEnabledNodeGrant(this.AdminContext(), &pb.FindEnabledNodeGrantRequest{NodeGrantId: params.GrantId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	if grantResp.NodeGrant == nil {
		this.WriteString("can not find the grant")
		return
	}

	grant := grantResp.NodeGrant
	this.Data["grant"] = maps.Map{
		"id":          grant.Id,
		"name":        grant.Name,
		"method":      grant.Method,
		"methodName":  grantutils.FindGrantMethodName(grant.Method, this.LangCode()),
		"username":    grant.Username,
		"password":    strings.Repeat("*", len(grant.Password)),
		"privateKey":  grant.PrivateKey,
		"description": grant.Description,
		"su":          grant.Su,
	}

	this.Show()
}

func (this *TestAction) RunPost(params struct {
	GrantId int64
	Host    string
	Port    int32

	Must *actions.Must
	CSRF *actionutils.CSRF
}) {
	params.Must.
		Field("host", params.Host).
		Require("请输入节点主机地址").
		Field("port", params.Port).
		Gt(0, "请输入正确的端口号").
		Lt(65535, "请输入正确的端口号")

	resp, err := this.RPC().NodeGrantRPC().TestNodeGrant(this.AdminContext(), &pb.TestNodeGrantRequest{
		NodeGrantId: params.GrantId,
		Host:        params.Host,
		Port:        params.Port,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	this.Data["isOk"] = resp.IsOk
	this.Data["error"] = resp.Error
	this.Success()
}
