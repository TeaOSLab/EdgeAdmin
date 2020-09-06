package api

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/rpc/pb"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/iwind/TeaGo/maps"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.Nav("", "node", "index")
}

func (this *IndexAction) RunGet(params struct{}) {
	countResp, err := this.RPC().APINodeRPC().CountAllEnabledAPINodes(this.AdminContext(), &pb.CountAllEnabledAPINodesRequest{})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	count := countResp.Count
	page := this.NewPage(count)
	this.Data["page"] = page.AsHTML()

	nodeMaps := []maps.Map{}
	if count > 0 {
		nodesResp, err := this.RPC().APINodeRPC().ListEnabledAPINodes(this.AdminContext(), &pb.ListEnabledAPINodesRequest{
			Offset: page.Offset,
			Size:   page.Size,
		})
		if err != nil {
			this.ErrorPage(err)
			return
		}

		for _, node := range nodesResp.Nodes {
			nodeMaps = append(nodeMaps, maps.Map{
				"id":   node.Id,
				"isOn": node.IsOn,
				"name": node.Name,
				"host": node.Host,
				"port": node.Port,
			})
		}
	}
	this.Data["nodes"] = nodeMaps

	this.Show()
}
