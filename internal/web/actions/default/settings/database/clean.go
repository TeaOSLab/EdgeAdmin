package profile

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/utils/numberutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
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

func (this *CleanAction) RunGet(params struct{}) {
	this.Show()
}

func (this *CleanAction) RunPost(params struct {
	Must *actions.Must
}) {
	tablesResp, err := this.RPC().DBRPC().FindAllDBTables(this.AdminContext(), &pb.FindAllDBTablesRequest{})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	tableMaps := []maps.Map{}
	for _, table := range tablesResp.DbTables {
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
