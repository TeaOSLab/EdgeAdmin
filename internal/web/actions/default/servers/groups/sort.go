package groups

import (	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
)

type SortAction struct {
	actionutils.ParentAction
}

func (this *SortAction) RunPost(params struct {
	GroupIds []int64
}) {
	// 创建日志
	defer this.CreateLogInfo(codes.ServerGroup_LogSortServerGroups)

	_, err := this.RPC().ServerGroupRPC().UpdateServerGroupOrders(this.AdminContext(), &pb.UpdateServerGroupOrdersRequest{ServerGroupIds: params.GroupIds})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
