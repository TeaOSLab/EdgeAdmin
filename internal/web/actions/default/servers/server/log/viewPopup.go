package log

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/iplibrary"
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
					"id":       policyResp.HttpFirewallPolicy.Id,
					"name":     policyResp.HttpFirewallPolicy.Name,
					"serverId": policyResp.HttpFirewallPolicy.ServerId,
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
	var ipRegion = iplibrary.LookupIP(accessLog.RemoteAddr)
	if ipRegion != nil && ipRegion.IsOk() {
		regionMap = maps.Map{
			"full": ipRegion.RegionSummary(),
			"isp":  ipRegion.ProviderName(),
		}
	}
	this.Data["region"] = regionMap

	// 请求内容
	this.Data["requestBody"] = string(accessLog.RequestBody)
	this.Data["requestContentType"] = "text/plain"

	requestContentType, ok := accessLog.Header["Content-Type"]
	if ok {
		if len(requestContentType.Values) > 0 {
			var contentType = requestContentType.Values[0]
			if strings.HasSuffix(contentType, "/json") || strings.Contains(contentType, "/json;") {
				this.Data["requestContentType"] = "application/json"
			}
		}
	}

	this.Show()
}
