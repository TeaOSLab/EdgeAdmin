package usernodes

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/maps"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.Nav("", "node", "index")
}

func (this *IndexAction) RunGet(params struct{}) {
	countResp, err := this.RPC().UserNodeRPC().CountAllEnabledUserNodes(this.AdminContext(), &pb.CountAllEnabledUserNodesRequest{})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	count := countResp.Count
	page := this.NewPage(count)
	this.Data["page"] = page.AsHTML()

	nodeMaps := []maps.Map{}
	if count > 0 {
		nodesResp, err := this.RPC().UserNodeRPC().ListEnabledUserNodes(this.AdminContext(), &pb.ListEnabledUserNodesRequest{
			Offset: page.Offset,
			Size:   page.Size,
		})
		if err != nil {
			this.ErrorPage(err)
			return
		}

		for _, node := range nodesResp.Nodes {
			nodeMaps = append(nodeMaps, maps.Map{
				"id":          node.Id,
				"isOn":        node.IsOn,
				"name":        node.Name,
				"accessAddrs": node.AccessAddrs,
			})
		}
	}
	this.Data["nodes"] = nodeMaps

	this.Show()
}
