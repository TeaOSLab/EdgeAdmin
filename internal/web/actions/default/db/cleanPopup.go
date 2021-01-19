package db

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/utils/numberutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type CleanPopupAction struct {
	actionutils.ParentAction
}

func (this *CleanPopupAction) Init() {
	this.Nav("", "", "")
}

func (this *CleanPopupAction) RunGet(params struct {
	NodeId int64
}) {
	this.Data["nodeId"] = params.NodeId

	this.Show()
}

func (this *CleanPopupAction) RunPost(params struct {
	NodeId int64

	Must *actions.Must
}) {
	tablesResp, err := this.RPC().DBNodeRPC().FindAllDBNodeTables(this.AdminContext(), &pb.FindAllDBNodeTablesRequest{
		DbNodeId: params.NodeId,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	tableMaps := []maps.Map{}
	for _, table := range tablesResp.DbNodeTables {
		if !table.IsBaseTable || (!table.CanClean && !table.CanDelete) {
			continue
		}
		tableMaps = append(tableMaps, maps.Map{
			"name":      table.Name,
			"rows":      table.Rows,
			"size":      numberutils.FormatBytes(table.DataLength + table.IndexLength),
			"canDelete": table.CanDelete,
			"canClean":  table.CanClean,
			"comment":   table.Comment,
		})
	}
	this.Data["tables"] = tableMaps
	this.Success()
}
