package ipadmin

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
)

type DeleteIPAction struct {
	actionutils.ParentAction
}

func (this *DeleteIPAction) RunPost(params struct {
	ItemId int64
}) {
	// TODO 判断权限

	_, err := this.RPC().IPItemRPC().DeleteIPItem(this.AdminContext(), &pb.DeleteIPItemRequest{IpItemId: params.ItemId})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
