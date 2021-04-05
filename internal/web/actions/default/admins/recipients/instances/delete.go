package instances

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
)

type DeleteAction struct {
	actionutils.ParentAction
}

func (this *DeleteAction) RunPost(params struct {
	InstanceId int64
}) {
	defer this.CreateLogInfo("删除消息媒介 %d", params.InstanceId)

	_, err := this.RPC().MessageMediaInstanceRPC().DeleteMessageMediaInstance(this.AdminContext(), &pb.DeleteMessageMediaInstanceRequest{MessageMediaInstanceId: params.InstanceId})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
