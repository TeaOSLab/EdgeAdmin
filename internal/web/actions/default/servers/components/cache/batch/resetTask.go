// Copyright 2022 GoEdge CDN goedge.cdn@gmail.com. All rights reserved. Official site: https://goedge.cn .

package cache

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
)

type ResetTaskAction struct {
	actionutils.ParentAction
}

func (this *ResetTaskAction) RunPost(params struct {
	TaskId int64
}) {
	this.CreateLogInfo(codes.HTTPCacheTask_LogResetHTTPCacheTask, params.TaskId)

	_, err := this.RPC().HTTPCacheTaskRPC().ResetHTTPCacheTask(this.AdminContext(), &pb.ResetHTTPCacheTaskRequest{HttpCacheTaskId: params.TaskId})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
