package database

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/utils/numberutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
	"sort"
)

type CleanAction struct {
	actionutils.ParentAction
}

func (this *CleanAction) Init() {
	this.Nav("", "", "clean")
}

func (this *CleanAction) RunGet(params struct {
	OrderTable string
	OrderSize  string
}) {
	this.Data["orderTable"] = params.OrderTable
	this.Data["orderSize"] = params.OrderSize

	this.Show()
}

func (this *CleanAction) RunPost(params struct {
	OrderTable string
	OrderSize  string

	Must *actions.Must
}) {
	tablesResp, err := this.RPC().DBRPC().FindAllDBTables(this.AdminContext(), &pb.FindAllDBTablesRequest{})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	var tables = tablesResp.DbTables

	// 排序
	switch params.OrderTable {
	case "asc":
		sort.Slice(tables, func(i, j int) bool {
			return tables[i].Name < tables[j].Name
		})
	case "desc":
		sort.Slice(tables, func(i, j int) bool {
			return tables[i].Name > tables[j].Name
		})
	}

	switch params.OrderSize {
	case "asc":
		sort.Slice(tables, func(i, j int) bool {
			return tables[i].DataLength+tables[i].IndexLength < tables[j].DataLength+tables[j].IndexLength
		})
	case "desc":
		sort.Slice(tables, func(i, j int) bool {
			return tables[i].DataLength+tables[i].IndexLength > tables[j].DataLength+tables[j].IndexLength
		})
	}

	var tableMaps = []maps.Map{}
	for _, table := range tables {
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
