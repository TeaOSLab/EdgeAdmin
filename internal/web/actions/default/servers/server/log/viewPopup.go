package log

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/maps"
	"net/http"
	"strings"
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
	accessLog := accessLogResp.AccessLog
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
		policyResp, err := this.RPC().HTTPFirewallPolicyRPC().FindEnabledFirewallPolicy(this.AdminContext(), &pb.FindEnabledFirewallPolicyRequest{FirewallPolicyId: accessLog.FirewallPolicyId})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		if policyResp.FirewallPolicy != nil {
			wafMap = maps.Map{
				"policy": maps.Map{
					"id":   policyResp.FirewallPolicy.Id,
					"name": policyResp.FirewallPolicy.Name,
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
	region := regionResp.Region
	if region != nil {
		pieces := []string{}
		if len(region.Country) > 0 {
			pieces = append(pieces, region.Country)
		}
		if len(region.Region) > 0 {
			pieces = append(pieces, region.Region)
		}
		if len(region.Province) > 0 {
			pieces = append(pieces, region.Province)
		}
		if len(region.City) > 0 {
			pieces = append(pieces, region.City)
		}
		regionMap = maps.Map{
			"full": strings.Join(pieces, " "),
			"isp":  region.Isp,
		}
	}
	this.Data["region"] = regionMap

	this.Show()
}
