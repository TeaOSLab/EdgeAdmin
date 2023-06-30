package instances

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
)

type DeleteAction struct {
	actionutils.ParentAction
}

func (this *DeleteAction) RunPost(params struct {
	InstanceId int64
}) {
	defer this.CreateLogInfo(codes.MessageMediaInstance_LogDeleteMessageMediaInstance, params.InstanceId)

	_, err := this.RPC().MessageMediaInstanceRPC().DeleteMessageMediaInstance(this.AdminContext(), &pb.DeleteMessageMediaInstanceRequest{MessageMediaInstanceId: params.InstanceId})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
