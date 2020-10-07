package waf

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
)

type UpdateGroupOnAction struct {
	actionutils.ParentAction
}

func (this *UpdateGroupOnAction) RunPost(params struct {
	GroupId int64
	IsOn    bool
}) {
	_, err := this.RPC().HTTPFirewallRuleGroupRPC().UpdateHTTPFirewallRuleGroupIsOn(this.AdminContext(), &pb.UpdateHTTPFirewallRuleGroupIsOnRequest{
		FirewallRuleGroupId: params.GroupId,
		IsOn:                params.IsOn,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
