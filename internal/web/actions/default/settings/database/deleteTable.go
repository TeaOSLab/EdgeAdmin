package profile

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
)

type DeleteTableAction struct {
	actionutils.ParentAction
}

func (this *DeleteTableAction) RunPost(params struct {
	Table string
}) {
	defer this.CreateLogInfo("删除数据表 %s", params.Table)

	_, err := this.RPC().DBRPC().DeleteDBTable(this.AdminContext(), &pb.DeleteDBTableRequest{DbTable: params.Table})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	this.Success()
}
