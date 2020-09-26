package clusters

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/configutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/types"
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
			// 全部节点数量
			countNodesResp, err := this.RPC().NodeRPC().CountAllEnabledNodesMatch(this.AdminContext(), &pb.CountAllEnabledNodesMatchRequest{ClusterId: cluster.Id})
			if err != nil {
				this.ErrorPage(err)
				return
			}

			// 在线节点
			countActiveNodesResp, err := this.RPC().NodeRPC().CountAllEnabledNodesMatch(this.AdminContext(), &pb.CountAllEnabledNodesMatchRequest{
				ClusterId:   cluster.Id,
				ActiveState: types.Int32(configutils.BoolStateYes),
			})
			if err != nil {
				this.ErrorPage(err)
				return
			}

			clusterMaps = append(clusterMaps, maps.Map{
				"id":               cluster.Id,
				"name":             cluster.Name,
				"installDir":       cluster.InstallDir,
				"hasGrant":         cluster.GrantId > 0,
				"countAllNodes":    countNodesResp.Count,
				"countActiveNodes": countActiveNodesResp.Count,
			})
		}
	}
	this.Data["clusters"] = clusterMaps

	this.Show()
}
