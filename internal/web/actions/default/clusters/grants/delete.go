package grants

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
)

type DeleteAction struct {
	actionutils.ParentAction
}

func (this *DeleteAction) RunPost(params struct {
	GrantId int64
}) {
	// 检查是否有别的集群或节点正在使用
	countResp, err := this.RPC().NodeClusterRPC().CountAllEnabledNodeClustersWithGrantId(this.AdminContext(), &pb.CountAllEnabledNodeClustersWithGrantIdRequest{
		GrantId: params.GrantId,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	if countResp.Count > 0 {
		this.Fail("有集群正在使用此服务，请修改后再删除")
	}

	countResp2, err := this.RPC().NodeRPC().CountAllEnabledNodesWithGrantId(this.AdminContext(), &pb.CountAllEnabledNodesWithGrantIdRequest{GrantId: params.GrantId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	if countResp2.Count > 0 {
		this.Fail("有节点正在使用此服务，请修改后再删除")
	}

	// 删除
	_, err = this.RPC().NodeGrantRPC().DisableNodeGrant(this.AdminContext(), &pb.DisableNodeGrantRequest{GrantId: params.GrantId})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
