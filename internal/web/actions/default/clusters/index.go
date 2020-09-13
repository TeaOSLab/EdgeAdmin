package clusters

import (
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/iwind/TeaGo/maps"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.Nav("", "cluster", "index")
}

func (this *IndexAction) RunGet(params struct{}) {
	countResp, err := this.RPC().NodeClusterRPC().CountAllEnabledNodeClusters(this.AdminContext(), &pb.CountAllEnabledNodeClustersRequest{})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	count := countResp.Count
	page := this.NewPage(count)
	this.Data["page"] = page.AsHTML()

	clusterMaps := []maps.Map{}
	if count > 0 {
		clustersResp, err := this.RPC().NodeClusterRPC().ListEnabledNodeClusters(this.AdminContext(), &pb.ListEnabledNodeClustersRequest{
			Offset: page.Offset,
			Size:   page.Size,
		})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		for _, cluster := range clustersResp.Clusters {
			// 节点数量
			countNodesResp, err := this.RPC().NodeRPC().CountAllEnabledNodesMatch(this.AdminContext(), &pb.CountAllEnabledNodesMatchRequest{ClusterId: cluster.Id})
			if err != nil {
				this.ErrorPage(err)
				return
			}

			clusterMaps = append(clusterMaps, maps.Map{
				"id":         cluster.Id,
				"name":       cluster.Name,
				"installDir": cluster.InstallDir,
				"hasGrant":   cluster.GrantId > 0,
				"countNodes": countNodesResp.Count,
			})
		}
	}
	this.Data["clusters"] = clusterMaps

	this.Show()
}
