// Copyright 2022 GoEdge CDN goedge.cdn@gmail.com. All rights reserved.

package ssh

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
)

type TestAction struct {
	actionutils.ParentAction
}

func (this *TestAction) RunPost(params struct {
	GrantId int64
	Host    string
	Port    int32
}) {
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
