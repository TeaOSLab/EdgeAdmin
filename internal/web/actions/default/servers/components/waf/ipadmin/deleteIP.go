package ipadmin

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/servers/components/waf/ipadmin/ipadminutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
)

type DeleteIPAction struct {
	actionutils.ParentAction
}

func (this *DeleteIPAction) RunPost(params struct {
	FirewallPolicyId int64
	ItemId           int64
}) {
	// TODO 判断权限

	_, err := this.RPC().IPItemRPC().DeleteIPItem(this.AdminContext(), &pb.DeleteIPItemRequest{IpItemId: params.ItemId})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	// 发送通知
	err = ipadminutils.NotifyUpdateToClustersWithFirewallPolicyId(this.AdminContext(), params.FirewallPolicyId)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
