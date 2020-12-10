package items

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
)

type DeleteAction struct {
	actionutils.ParentAction
}

func (this *DeleteAction) RunPost(params struct {
	ItemId int64
}) {
	defer this.CreateLogInfo("删除流量价格项目 %d", params.ItemId)

	_, err := this.RPC().NodePriceItemRPC().DeleteNodePriceItem(this.AdminContext(), &pb.DeleteNodePriceItemRequest{NodePriceItemId: params.ItemId})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
