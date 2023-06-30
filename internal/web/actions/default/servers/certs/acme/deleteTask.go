package acme

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
)

type DeleteTaskAction struct {
	actionutils.ParentAction
}

func (this *DeleteTaskAction) RunPost(params struct {
	TaskId int64
}) {
	defer this.CreateLogInfo(codes.ACMETask_LogDeleteACMETask, params.TaskId)

	_, err := this.RPC().ACMETaskRPC().DeleteACMETask(this.AdminContext(), &pb.DeleteACMETaskRequest{AcmeTaskId: params.TaskId})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
