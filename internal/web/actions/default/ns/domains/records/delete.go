// Copyright 2021 Liuxiangchao iwind.liu@gmail.com. All rights reserved.

package records

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
)

type DeleteAction struct {
	actionutils.ParentAction
}

func (this *DeleteAction) RunPost(params struct {
	RecordId int64
}) {
	defer this.CreateLogInfo("删除域名记录 %d", params.RecordId)

	_, err := this.RPC().NSRecordRPC().DeleteNSRecord(this.AdminContext(), &pb.DeleteNSRecordRequest{NsRecordId: params.RecordId})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
