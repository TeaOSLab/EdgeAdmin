package waf

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/oplogs"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
)

type DeleteAction struct {
	actionutils.ParentAction
}

func (this *DeleteAction) RunPost(params struct {
	FirewallPolicyId int64
}) {
	// 日志
	this.CreateLog(oplogs.LevelInfo, "删除WAF策略 %d", params.FirewallPolicyId)

	countResp, err := this.RPC().ServerRPC().CountAllEnabledServersWithHTTPFirewallPolicyId(this.AdminContext(), &pb.CountAllEnabledServersWithHTTPFirewallPolicyIdRequest{FirewallPolicyId: params.FirewallPolicyId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	if countResp.Count > 0 {
		this.Fail("此WAF策略正在被有些服务引用，请修改后再删除。")
	}

	_, err = this.RPC().HTTPFirewallPolicyRPC().DeleteFirewallPolicy(this.AdminContext(), &pb.DeleteFirewallPolicyRequest{FirewallPolicyId: params.FirewallPolicyId})
	if err != nil {
		this.ErrorPage(err)
		return
	}



	this.Success()
}
