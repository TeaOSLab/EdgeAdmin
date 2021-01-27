package tasks

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
)

type CheckAction struct {
	actionutils.ParentAction
}

func (this *CheckAction) RunPost(params struct{}) {
	resp, err := this.RPC().DNSTaskRPC().ExistsDNSTasks(this.AdminContext(), &pb.ExistsDNSTasksRequest{})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Data["isDoing"] = resp.ExistTasks
	this.Data["hasError"] = resp.ExistError

	this.Success()
}
