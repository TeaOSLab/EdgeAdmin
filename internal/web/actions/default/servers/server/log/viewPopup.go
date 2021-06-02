package log

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/maps"
	"net/http"
)

type ViewPopupAction struct {
	actionutils.ParentAction
}

func (this *ViewPopupAction) Init() {
	this.Nav("", "", "")
}

func (this *ViewPopupAction) RunGet(params struct {
	RequestId string
}) {
	accessLogResp, err := this.RPC().HTTPAccessLogRPC().FindHTTPAccessLog(this.AdminContext(), &pb.FindHTTPAccessLogRequest{RequestId: params.RequestId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	accessLog := accessLogResp.HttpAccessLog
	if accessLog == nil {
		this.WriteString("not found: " + params.RequestId)
		return
	}

	// 状态
	if len(accessLog.StatusMessage) == 0 {
		accessLog.StatusMessage = http.StatusText(int(accessLog.Status))
	}

	this.Data["accessLog"] = accessLog

	// WAF相关
	var wafMap maps.Map = nil
	if accessLog.FirewallPolicyId > 0 {
		policyResp, err := this.RPC().HTTPFirewallPolicyRPC().FindEnabledHTTPFirewallPolicy(this.AdminContext(), &pb.FindEnabledHTTPFirewallPolicyRequest{HttpFirewallPolicyId: accessLog.FirewallPolicyId})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		if policyResp.HttpFirewallPolicy != nil {
			wafMap = maps.Map{
				"policy": maps.Map{
					"id":   policyResp.HttpFirewallPolicy.Id,
					"name": policyResp.HttpFirewallPolicy.Name,
				},
			}
			if accessLog.FirewallRuleGroupId > 0 {
				groupResp, err := this.RPC().HTTPFirewallRuleGroupRPC().FindEnabledHTTPFirewallRuleGroup(this.AdminContext(), &pb.FindEnabledHTTPFirewallRuleGroupRequest{FirewallRuleGroupId: accessLog.FirewallRuleGroupId})
				if err != nil {
					this.ErrorPage(err)
					return
				}
				if groupResp.FirewallRuleGroup != nil {
					wafMap["group"] = maps.Map{
						"id":   groupResp.FirewallRuleGroup.Id,
						"name": groupResp.FirewallRuleGroup.Name,
					}

					if accessLog.FirewallRuleSetId > 0 {
						setResp, err := this.RPC().HTTPFirewallRuleSetRPC().FindEnabledHTTPFirewallRuleSet(this.AdminContext(), &pb.FindEnabledHTTPFirewallRuleSetRequest{FirewallRuleSetId: accessLog.FirewallRuleSetId})
						if err != nil {
							this.ErrorPage(err)
							return
						}
						if setResp.FirewallRuleSet != nil {
							wafMap["set"] = maps.Map{
								"id":   setResp.FirewallRuleSet.Id,
								"name": setResp.FirewallRuleSet.Name,
							}
						}
					}
				}
			}
		}
	}
	this.Data["wafInfo"] = wafMap

	// 地域相关
	var regionMap maps.Map = nil
	regionResp, err := this.RPC().IPLibraryRPC().LookupIPRegion(this.AdminContext(), &pb.LookupIPRegionRequest{Ip: accessLog.RemoteAddr})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	region := regionResp.IpRegion
	if region != nil {
		regionMap = maps.Map{
			"full": region.Summary,
			"isp":  region.Isp,
		}
	}
	this.Data["region"] = regionMap

	this.Show()
}
