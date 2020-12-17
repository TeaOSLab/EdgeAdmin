package clusters

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/maps"
)

// 检查变更的集群列表
type CheckChangeAction struct {
	actionutils.ParentAction
}

func (this *CheckChangeAction) Init() {
	this.Nav("", "", "")
}

func (this *CheckChangeAction) RunPost(params struct {
	IsNotifying bool
}) {
	resp, err := this.RPC().NodeClusterRPC().FindAllChangedNodeClusters(this.AdminContext(), &pb.FindAllChangedNodeClustersRequest{})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	result := []maps.Map{}
	for _, cluster := range resp.NodeClusters {
		result = append(result, maps.Map{
			"id":   cluster.Id,
			"name": cluster.Name,
		})
	}

	this.Data["clusters"] = result
	this.Success()
}
