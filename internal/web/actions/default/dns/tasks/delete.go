package tasks

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
)

type DeleteAction struct {
	actionutils.ParentAction
}

func (this *DeleteAction) RunPost(params struct {
	TaskId int64
}) {
	defer this.CreateLogInfo(codes.DNSTask_LogDeleteDNSTask, params.TaskId)

	_, err := this.RPC().DNSTaskRPC().DeleteDNSTask(this.AdminContext(), &pb.DeleteDNSTaskRequest{DnsTaskId: params.TaskId})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
