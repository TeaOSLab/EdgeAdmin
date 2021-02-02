package tasks

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
)

type DeleteBatchAction struct {
	actionutils.ParentAction
}

func (this *DeleteBatchAction) RunPost(params struct {
	TaskIds []int64
}) {
	defer this.CreateLogInfo("批量删除节点同步任务")

	_, err := this.RPC().NodeTaskRPC().DeleteNodeTasks(this.AdminContext(), &pb.DeleteNodeTasksRequest{NodeTaskIds: params.TaskIds})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
