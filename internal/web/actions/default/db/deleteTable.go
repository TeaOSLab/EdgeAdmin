package db

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
)

type DeleteTableAction struct {
	actionutils.ParentAction
}

func (this *DeleteTableAction) RunPost(params struct {
	NodeId int64
	Table  string
}) {
	defer this.CreateLogInfo(codes.DBNode_LogDeleteTable, params.NodeId, params.Table)

	_, err := this.RPC().DBNodeRPC().DeleteDBNodeTable(this.AdminContext(), &pb.DeleteDBNodeTableRequest{
		DbNodeId:    params.NodeId,
		DbNodeTable: params.Table,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	this.Success()
}
