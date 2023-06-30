package database

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
)

type DeleteTableAction struct {
	actionutils.ParentAction
}

func (this *DeleteTableAction) RunPost(params struct {
	Table string
}) {
	defer this.CreateLogInfo(codes.Database_LogDeleteTable, params.Table)

	_, err := this.RPC().DBRPC().DeleteDBTable(this.AdminContext(), &pb.DeleteDBTableRequest{DbTable: params.Table})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	this.Success()
}
