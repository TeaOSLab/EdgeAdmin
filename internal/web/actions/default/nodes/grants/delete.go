package grants

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/rpc/pb"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
)

type DeleteAction struct {
	actionutils.ParentAction
}

func (this *DeleteAction) RunPost(params struct {
	GrantId int64
}) {
	_, err := this.RPC().NodeGrantRPC().DisableNodeGrant(this.AdminContext(), &pb.DisableNodeGrantRequest{GrantId: params.GrantId})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
