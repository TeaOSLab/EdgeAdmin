package clusters

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/maps"
)

type OptionsAction struct {
	actionutils.ParentAction
}

func (this *OptionsAction) RunPost(params struct{}) {
	clustersResp, err := this.RPC().NSClusterRPC().FindAllEnabledNSClusters(this.AdminContext(), &pb.FindAllEnabledNSClustersRequest{})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	clusterMaps := []maps.Map{}
	for _, cluster := range clustersResp.NsClusters {
		clusterMaps = append(clusterMaps, maps.Map{
			"id":   cluster.Id,
			"name": cluster.Name,
		})
	}
	this.Data["clusters"] = clusterMaps

	this.Success()
}
