package waf

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/dao"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.Nav("", "setting", "")
	this.SecondMenu("waf")
}

func (this *IndexAction) RunGet(params struct {
	ClusterId int64
}) {
	cluster, err := dao.SharedNodeClusterDAO.FindEnabledNodeCluster(this.AdminContext(), params.ClusterId)
	if err != nil {
		this.ErrorPage(err)
		return
	}
	if cluster == nil {
		this.NotFound("nodeCluster", params.ClusterId)
		return
	}

	// WAF设置
	this.Data["firewallPolicy"] = nil
	if cluster.HttpFirewallPolicyId > 0 {
		firewallPolicy, err := dao.SharedHTTPFirewallPolicyDAO.FindEnabledHTTPFirewallPolicy(this.AdminContext(), cluster.HttpFirewallPolicyId)
		if err != nil {
			this.ErrorPage(err)
			return
		}
		if firewallPolicy != nil {
			this.Data["firewallPolicy"] = maps.Map{
				"id":   firewallPolicy.Id,
				"name": firewallPolicy.Name,
				"isOn": firewallPolicy.IsOn,
			}
		}
	}

	this.Show()
}

func (this *IndexAction) RunPost(params struct {
	ClusterId            int64
	HttpFirewallPolicyId int64

	Must *actions.Must
	CSRF *actionutils.CSRF
}) {
	defer this.CreateLogInfo(codes.WAFPolicy_LogUpdateClusterWAFPolicy, params.ClusterId, params.HttpFirewallPolicyId)

	if params.HttpFirewallPolicyId <= 0 {
		this.Fail("请选择WAF策略")
	}

	_, err := this.RPC().NodeClusterRPC().UpdateNodeClusterHTTPFirewallPolicyId(this.AdminContext(), &pb.UpdateNodeClusterHTTPFirewallPolicyIdRequest{
		NodeClusterId:        params.ClusterId,
		HttpFirewallPolicyId: params.HttpFirewallPolicyId,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	this.Success()
}
