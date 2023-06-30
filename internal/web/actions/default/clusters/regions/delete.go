package regions

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
)

type DeleteAction struct {
	actionutils.ParentAction
}

func (this *DeleteAction) RunPost(params struct {
	RegionId int64
}) {
	defer this.CreateLogInfo(codes.NodeRegion_LogDeleteNodeRegion, params.RegionId)

	// 检查有无在使用
	countResp, err := this.RPC().NodeRPC().CountAllEnabledNodesWithNodeRegionId(this.AdminContext(), &pb.CountAllEnabledNodesWithNodeRegionIdRequest{NodeRegionId: params.RegionId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	if countResp.Count > 0 {
		this.Fail("此区域正在使用，不能删除")
	}

	// 执行删除
	_, err = this.RPC().NodeRegionRPC().DeleteNodeRegion(this.AdminContext(), &pb.DeleteNodeRegionRequest{NodeRegionId: params.RegionId})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
