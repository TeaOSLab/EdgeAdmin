package delete

import (
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/iwind/TeaGo/actions"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.Nav("", "delete", "")
	this.SecondMenu("index")
}

func (this *IndexAction) RunGet(params struct{}) {
	this.Show()
}

func (this *IndexAction) RunPost(params struct {
	ServerId int64
	Must     *actions.Must
}) {
	_, err := this.RPC().ServerRPC().DisableServer(this.AdminContext(), &pb.DisableServerRequest{ServerId: params.ServerId})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
