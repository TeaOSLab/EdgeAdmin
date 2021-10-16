// Copyright 2021 Liuxiangchao iwind.liu@gmail.com. All rights reserved.

package dns

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
)

type RegenerateCNAMEAction struct {
	actionutils.ParentAction
}

func (this *RegenerateCNAMEAction) RunPost(params struct {
	ServerId int64
}) {
	defer this.CreateLogInfo("重新生成服务 %d 的CNAME", params.ServerId)

	_, err := this.RPC().ServerRPC().RegenerateServerCNAME(this.AdminContext(), &pb.RegenerateServerCNAMERequest{ServerId: params.ServerId})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
