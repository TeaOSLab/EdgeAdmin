package profile

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
)

type TruncateTableAction struct {
	actionutils.ParentAction
}

func (this *TruncateTableAction) RunPost(params struct {
	Table string
}) {
	defer this.CreateLogInfo("清空数据表 %s 数据", params.Table)

	_, err := this.RPC().DBRPC().TruncateDBTable(this.AdminContext(), &pb.TruncateDBTableRequest{DbTable: params.Table})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	this.Success()
}
