package messages

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/oplogs"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
)

type ReadAllAction struct {
	actionutils.ParentAction
}

func (this *ReadAllAction) RunPost(params struct{}) {
	// 创建日志
	defer this.CreateLog(oplogs.LevelInfo, "将所有消息置为已读")

	_, err := this.RPC().MessageRPC().UpdateAllMessagesRead(this.AdminContext(), &pb.UpdateAllMessagesReadRequest{})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
