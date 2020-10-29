package groups

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
)

type SortAction struct {
	actionutils.ParentAction
}

func (this *SortAction) RunPost(params struct {
	GroupIds []int64
}) {
	_, err := this.RPC().ServerGroupRPC().UpdateServerGroupOrders(this.AdminContext(), &pb.UpdateServerGroupOrdersRequest{GroupIds: params.GroupIds})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
