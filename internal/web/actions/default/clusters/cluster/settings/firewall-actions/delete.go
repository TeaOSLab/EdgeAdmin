package firewallActions

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
)

type DeleteAction struct {
	actionutils.ParentAction
}

func (this *DeleteAction) RunPost(params struct {
	ActionId int64
}) {
	defer this.CreateLogInfo(codes.WAFAction_LogDeleteWAFAction, params.ActionId)

	_, err := this.RPC().NodeClusterFirewallActionRPC().DeleteNodeClusterFirewallAction(this.AdminContext(), &pb.DeleteNodeClusterFirewallActionRequest{NodeClusterFirewallActionId: params.ActionId})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
