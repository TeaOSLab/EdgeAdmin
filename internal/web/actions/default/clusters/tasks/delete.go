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
	defer this.CreateLogInfo(codes.NodeTask_LogDeleteNodeTask, params.TaskId)

	_, err := this.RPC().NodeTaskRPC().DeleteNodeTask(this.AdminContext(), &pb.DeleteNodeTaskRequest{NodeTaskId: params.TaskId})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
