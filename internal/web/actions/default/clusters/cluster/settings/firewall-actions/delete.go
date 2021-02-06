package firewallActions

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
)

type DeleteAction struct {
	actionutils.ParentAction
}

func (this *DeleteAction) RunPost(params struct {
	ActionId int64
}) {
	defer this.CreateLogInfo("删除WAF动作 %d", params.ActionId)

	_, err := this.RPC().NodeClusterFirewallActionRPC().DeleteNodeClusterFirewallAction(this.AdminContext(), &pb.DeleteNodeClusterFirewallActionRequest{NodeClusterFirewallActionId: params.ActionId})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
