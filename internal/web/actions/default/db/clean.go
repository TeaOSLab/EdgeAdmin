package db

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/utils/numberutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/db/dbnodeutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type CleanAction struct {
	actionutils.ParentAction
}

func (this *CleanAction) Init() {
	this.Nav("", "", "clean")
}

func (this *CleanAction) RunGet(params struct {
	NodeId int64
}) {
	_, err := dbnodeutils.InitNode(this.Parent(), params.NodeId)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Data["nodeId"] = params.NodeId

	this.Show()
}

func (this *CleanAction) RunPost(params struct {
	NodeId int64

	Must *actions.Must
}) {
	tablesResp, err := this.RPC().DBNodeRPC().FindAllDBNodeTables(this.AdminContext(), &pb.FindAllDBNodeTablesRequest{
		DbNodeId: params.NodeId,
	})
	if err != nil {
		this.Fail("查询数据时出错了：" + err.Error())
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
