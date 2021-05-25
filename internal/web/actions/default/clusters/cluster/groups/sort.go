package groups

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/oplogs"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
)

type SortAction struct {
	actionutils.ParentAction
}

func (this *SortAction) RunPost(params struct {
	GroupIds []int64
}) {
	_, err := this.RPC().NodeGroupRPC().UpdateNodeGroupOrders(this.AdminContext(), &pb.UpdateNodeGroupOrdersRequest{NodeGroupIds: params.GroupIds})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	// 创建日志
	defer this.CreateLog(oplogs.LevelInfo, "修改集群分组排序")

	this.Success()
}
