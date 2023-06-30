package db

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
)

type TruncateTableAction struct {
	actionutils.ParentAction
}

func (this *TruncateTableAction) RunPost(params struct {
	NodeId int64
	Table  string
}) {
	defer this.CreateLogInfo(codes.DBNode_LogTruncateTable, params.NodeId, params.Table)

	_, err := this.RPC().DBNodeRPC().TruncateDBNodeTable(this.AdminContext(), &pb.TruncateDBNodeTableRequest{
		DbNodeId:    params.NodeId,
		DbNodeTable: params.Table,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	this.Success()
}
