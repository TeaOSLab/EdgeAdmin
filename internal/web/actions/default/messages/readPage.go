package messages

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/oplogs"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
)

type ReadPageAction struct {
	actionutils.ParentAction
}

func (this *ReadPageAction) RunPost(params struct {
	MessageIds []int64
}) {
	// 创建日志
	defer this.CreateLog(oplogs.LevelInfo, "将一组消息置为已读")

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
