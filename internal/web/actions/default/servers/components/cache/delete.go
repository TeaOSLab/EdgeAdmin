package cache

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
)

type DeleteAction struct {
	actionutils.ParentAction
}

func (this *DeleteAction) RunPost(params struct {
	CachePolicyId int64
}) {
	// 检查是否被引用
	countResp, err := this.RPC().NodeClusterRPC().CountAllEnabledNodeClustersWithHTTPCachePolicyId(this.AdminContext(), &pb.CountAllEnabledNodeClustersWithHTTPCachePolicyIdRequest{HttpCachePolicyId: params.CachePolicyId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	if countResp.Count > 0 {
		this.Fail("此缓存策略正在被有些集群引用，请修改后再删除。")
	}

	_, err = this.RPC().HTTPCachePolicyRPC().DeleteHTTPCachePolicy(this.AdminContext(), &pb.DeleteHTTPCachePolicyRequest{HttpCachePolicyId: params.CachePolicyId})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	// 创建日志
	defer this.CreateLogInfo(codes.ServerCachePolicy_LogDeleteCachePolicy, params.CachePolicyId)

	this.Success()
}
