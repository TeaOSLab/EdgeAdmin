package tasks

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/maps"
)

type TaskInfoAction struct {
	actionutils.ParentAction
}

func (this *TaskInfoAction) Init() {
	this.Nav("", "", "")
}

func (this *TaskInfoAction) RunPost(params struct {
	TaskId int64
}) {
	resp, err := this.RPC().MessageTaskRPC().FindEnabledMessageTask(this.AdminContext(), &pb.FindEnabledMessageTaskRequest{MessageTaskId: params.TaskId})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	if resp.MessageTask == nil {
		this.NotFound("messageTask", params.TaskId)
		return
	}

	result := resp.MessageTask.Result
	this.Data["status"] = resp.MessageTask.Status
	this.Data["result"] = maps.Map{
		"isOk":     result.IsOk,
		"response": result.Response,
		"error":    result.Error,
	}

	this.Success()
}
