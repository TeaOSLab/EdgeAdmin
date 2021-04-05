package groups

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
)

type DeleteAction struct {
	actionutils.ParentAction
}

func (this *DeleteAction) RunPost(params struct {
	GroupId int64
}) {
	_, err := this.RPC().MessageRecipientGroupRPC().DeleteMessageRecipientGroup(this.AdminContext(), &pb.DeleteMessageRecipientGroupRequest{MessageRecipientGroupId: params.GroupId})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
