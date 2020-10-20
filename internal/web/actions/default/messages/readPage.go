package messages

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
)

type ReadPageAction struct {
	actionutils.ParentAction
}

func (this *ReadPageAction) RunPost(params struct {
	MessageIds []int64
}) {
	_, err := this.RPC().MessageRPC().UpdateMessagesRead(this.AdminContext(), &pb.UpdateMessagesReadRequest{
		MessageIds: params.MessageIds,
		IsRead:     true,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
