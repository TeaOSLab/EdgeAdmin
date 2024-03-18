package grants

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/clusters/grants/grantutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/maps"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.Nav("", "grant", "index")
}

func (this *IndexAction) RunGet(params struct {
	Keyword string
}) {
	this.Data["keyword"] = params.Keyword

	countResp, err := this.RPC().NodeGrantRPC().CountAllEnabledNodeGrants(this.AdminContext(), &pb.CountAllEnabledNodeGrantsRequest{
		Keyword: params.Keyword,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	var page = this.NewPage(countResp.Count)
	this.Data["page"] = page.AsHTML()

	grantsResp, err := this.RPC().NodeGrantRPC().ListEnabledNodeGrants(this.AdminContext(), &pb.ListEnabledNodeGrantsRequest{
		Keyword: params.Keyword,
		Offset:  page.Offset,
		Size:    page.Size,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	var grantMaps = []maps.Map{}
	for _, grant := range grantsResp.NodeGrants {
		// 集群数
		countClustersResp, err := this.RPC().NodeClusterRPC().CountAllEnabledNodeClustersWithNodeGrantId(this.AdminContext(), &pb.CountAllEnabledNodeClustersWithNodeGrantIdRequest{NodeGrantId: grant.Id})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		var countClusters = countClustersResp.Count

		// 节点数
		countNodesResp, err := this.RPC().NodeRPC().CountAllEnabledNodesWithNodeGrantId(this.AdminContext(), &pb.CountAllEnabledNodesWithNodeGrantIdRequest{NodeGrantId: grant.Id})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		var countNodes = countNodesResp.Count

		grantMaps = append(grantMaps, maps.Map{
			"id":   grant.Id,
			"name": grant.Name,
			"method": maps.Map{
				"type": grant.Method,
				"name": grantutils.FindGrantMethodName(grant.Method, this.LangCode()),
			},
			"username":      grant.Username,
			"countClusters": countClusters,
			"countNodes":    countNodes,
		})
	}
	this.Data["grants"] = grantMaps

	this.Show()
}
