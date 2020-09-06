package common

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/rpc/pb"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/iwind/TeaGo/maps"
)

type ChangedClustersAction struct {
	actionutils.ParentAction
}

func (this *ChangedClustersAction) Init() {
	this.Nav("", "", "")
}

func (this *ChangedClustersAction) RunGet(params struct{}) {
	resp, err := this.RPC().NodeClusterRPC().FindAllChangedNodeClusters(this.AdminContext(), &pb.FindAllChangedNodeClustersRequest{})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	result := []maps.Map{}
	for _, cluster := range resp.Clusters {
		result = append(result, maps.Map{
			"id":   cluster.Id,
			"name": cluster.Name,
		})
	}

	this.Data["clusters"] = result

	this.Success()
}
