package recipients

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
)

type DeleteAction struct {
	actionutils.ParentAction
}

func (this *DeleteAction) RunPost(params struct {
	RecipientId int64
}) {
	defer this.CreateLogInfo("删除媒介接收人 %d", params.RecipientId)

	_, err := this.RPC().MessageRecipientRPC().DeleteMessageRecipient(this.AdminContext(), &pb.DeleteMessageRecipientRequest{MessageRecipientId: params.RecipientId})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
