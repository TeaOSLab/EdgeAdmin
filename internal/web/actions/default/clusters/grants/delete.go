package grants

import (	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
)

type DeleteAction struct {
	actionutils.ParentAction
}

func (this *DeleteAction) RunPost(params struct {
	GrantId int64
}) {
	// 创建日志
	defer this.CreateLogInfo(codes.NodeGrant_LogDeleteSSHGrant, params.GrantId)

	// 检查是否有别的集群或节点正在使用
	countResp, err := this.RPC().NodeClusterRPC().CountAllEnabledNodeClustersWithNodeGrantId(this.AdminContext(), &pb.CountAllEnabledNodeClustersWithNodeGrantIdRequest{
		NodeGrantId: params.GrantId,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	if countResp.Count > 0 {
		this.Fail("有集群正在使用此服务，请修改后再删除")
	}

	countResp2, err := this.RPC().NodeRPC().CountAllEnabledNodesWithNodeGrantId(this.AdminContext(), &pb.CountAllEnabledNodesWithNodeGrantIdRequest{NodeGrantId: params.GrantId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	if countResp2.Count > 0 {
		this.Fail("有节点正在使用此服务，请修改后再删除")
	}

	// 删除
	_, err = this.RPC().NodeGrantRPC().DisableNodeGrant(this.AdminContext(), &pb.DisableNodeGrantRequest{NodeGrantId: params.GrantId})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
