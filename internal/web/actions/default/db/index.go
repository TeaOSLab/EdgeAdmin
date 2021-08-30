package db

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/utils/numberutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/maps"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.Nav("db", "db", "")
}

func (this *IndexAction) RunGet(params struct{}) {
	countResp, err := this.RPC().DBNodeRPC().CountAllEnabledDBNodes(this.AdminContext(), &pb.CountAllEnabledDBNodesRequest{})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	count := countResp.Count

	page := this.NewPage(count)
	listResp, err := this.RPC().DBNodeRPC().ListEnabledDBNodes(this.AdminContext(), &pb.ListEnabledDBNodesRequest{
		Offset: page.Offset,
		Size:   page.Size,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	nodeMaps := []maps.Map{}
	for _, node := range listResp.DbNodes {
		nodeMaps = append(nodeMaps, maps.Map{
			"id":       node.Id,
			"isOn":     node.IsOn,
			"name":     node.Name,
			"host":     node.Host,
			"port":     node.Port,
			"database": node.Database,
			"status": maps.Map{
				"isOk":    node.Status.IsOk,
				"error":   node.Status.Error,
				"size":    numberutils.FormatBytes(node.Status.Size),
				"version": node.Status.Version,
			},
		})
	}

	this.Data["nodes"] = nodeMaps
	this.Data["page"] = page.AsHTML()

	this.Show()
}
