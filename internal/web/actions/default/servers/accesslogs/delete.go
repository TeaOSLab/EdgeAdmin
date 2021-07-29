// Copyright 2021 Liuxiangchao iwind.liu@gmail.com. All rights reserved.

package ipbox

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
)

type DeleteAction struct {
	actionutils.ParentAction
}

func (this *DeleteAction) RunPost(params struct {
	PolicyId int64
}) {
	defer this.CreateLogInfo("删除访问日志策略 %d", params.PolicyId)

	_, err := this.RPC().HTTPAccessLogPolicyRPC().DeleteHTTPAccessLogPolicy(this.AdminContext(), &pb.DeleteHTTPAccessLogPolicyRequest{HttpAccessLogPolicyId: params.PolicyId})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
