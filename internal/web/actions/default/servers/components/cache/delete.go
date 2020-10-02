package cache

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
)

type DeleteAction struct {
	actionutils.ParentAction
}

func (this *DeleteAction) RunPost(params struct {
	CachePolicyId int64
}) {
	// 检查是否被引用
	countResp, err := this.RPC().ServerRPC().CountAllEnabledServersWithCachePolicyId(this.AdminContext(), &pb.CountAllEnabledServersWithCachePolicyIdRequest{CachePolicyId: params.CachePolicyId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	if countResp.Count > 0 {
		this.Fail("此缓存策略正在被别的服务引用，请修改后再删除。")
	}

	_, err = this.RPC().HTTPCachePolicyRPC().DeleteHTTPCachePolicy(this.AdminContext(), &pb.DeleteHTTPCachePolicyRequest{CachePolicyId: params.CachePolicyId})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
