package waf

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/iplibrary"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs/firewallconfigs"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/maps"
	timeutil "github.com/iwind/TeaGo/utils/time"
	"regexp"
	"strings"
)

type LogAction struct {
	actionutils.ParentAction
}

func (this *LogAction) Init() {
	this.Nav("", "", "log")
}

func (this *LogAction) RunGet(params struct {
	Day              string
	RequestId        string
	FirewallPolicyId int64
	GroupId          int64
	Partition        int32 `default:"-1"`
}) {
	if len(params.Day) == 0 {
		params.Day = timeutil.Format("Y-m-d")
	}

	this.Data["path"] = this.Request.URL.Path
	this.Data["day"] = params.Day
	this.Data["groupId"] = params.GroupId
	this.Data["accessLogs"] = []maps.Map{}
	this.Data["partition"] = params.Partition

	var day = params.Day
	var ipList = []string{}
	var wafMaps = []maps.Map{}
	if len(day) > 0 && regexp.MustCompile(`\d{4}-\d{2}-\d{2}`).MatchString(day) {
		day = strings.ReplaceAll(day, "-", "")
		var size = int64(20)

		resp, err := this.RPC().HTTPAccessLogRPC().ListHTTPAccessLogs(this.AdminContext(), &pb.ListHTTPAccessLogsRequest{
			Partition:           params.Partition,
			RequestId:           params.RequestId,
			FirewallPolicyId:    params.FirewallPolicyId,
			FirewallRuleGroupId: params.GroupId,
			Day:                 day,
			Size:                size,
		})
		if err != nil {
			this.ErrorPage(err)
			return
		}

		if len(resp.HttpAccessLogs) == 0 {
			this.Data["accessLogs"] = []interface{}{}
		} else {
			this.Data["accessLogs"] = resp.HttpAccessLogs
			for _, accessLog := range resp.HttpAccessLogs {
				// IP
				if len(accessLog.RemoteAddr) > 0 {
					if !lists.ContainsString(ipList, accessLog.RemoteAddr) {
						ipList = append(ipList, accessLog.RemoteAddr)
					}
				}

				// WAF信息集合
				if accessLog.FirewallPolicyId > 0 && accessLog.FirewallRuleGroupId > 0 && accessLog.FirewallRuleSetId > 0 {
					// 检查Set是否已经存在
					var existSet = false
					for _, wafMap := range wafMaps {
						if wafMap.GetInt64("setId") == accessLog.FirewallRuleSetId {
							existSet = true
							break
						}
					}
					if !existSet {
						wafMaps = append(wafMaps, maps.Map{
							"policyId": accessLog.FirewallPolicyId,
							"groupId":  accessLog.FirewallRuleGroupId,
							"setId":    accessLog.FirewallRuleSetId,
						})
					}
				}
			}
		}
		this.Data["hasMore"] = resp.HasMore
		this.Data["nextRequestId"] = resp.RequestId

		// 上一个requestId
		this.Data["hasPrev"] = false
		this.Data["lastRequestId"] = ""
		if len(params.RequestId) > 0 {
			this.Data["hasPrev"] = true
			prevResp, err := this.RPC().HTTPAccessLogRPC().ListHTTPAccessLogs(this.AdminContext(), &pb.ListHTTPAccessLogsRequest{
				Partition:           params.Partition,
				RequestId:           params.RequestId,
				FirewallPolicyId:    params.FirewallPolicyId,
				FirewallRuleGroupId: params.GroupId,
				Day:                 day,
				Size:                size,
				Reverse:             true,
			})
			if err != nil {
				this.ErrorPage(err)
				return
			}
			if int64(len(prevResp.HttpAccessLogs)) == size {
				this.Data["lastRequestId"] = prevResp.RequestId
			}
		}
	}

	// 所有分组
	policyResp, err := this.RPC().HTTPFirewallPolicyRPC().FindEnabledHTTPFirewallPolicyConfig(this.AdminContext(), &pb.FindEnabledHTTPFirewallPolicyConfigRequest{
		HttpFirewallPolicyId: params.FirewallPolicyId,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	policyConfig := &firewallconfigs.HTTPFirewallPolicy{}
	err = json.Unmarshal(policyResp.HttpFirewallPolicyJSON, policyConfig)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	groupMaps := []maps.Map{}
	for _, group := range policyConfig.AllRuleGroups() {
		groupMaps = append(groupMaps, maps.Map{
			"id":   group.Id,
			"name": group.Name,
		})
	}
	this.Data["groups"] = groupMaps

	// 根据IP查询区域
	this.Data["regions"] = iplibrary.LookupIPSummaries(ipList)

	// WAF相关
	var wafInfos = map[int64]maps.Map{}                          // set id => WAF Map
	var wafPolicyCacheMap = map[int64]*pb.HTTPFirewallPolicy{}   // id => *pb.HTTPFirewallPolicy
	var wafGroupCacheMap = map[int64]*pb.HTTPFirewallRuleGroup{} // id => *pb.HTTPFirewallRuleGroup
	var wafSetCacheMap = map[int64]*pb.HTTPFirewallRuleSet{}     // id => *pb.HTTPFirewallRuleSet
	for _, wafMap := range wafMaps {
		var policyId = wafMap.GetInt64("policyId")
		var groupId = wafMap.GetInt64("groupId")
		var setId = wafMap.GetInt64("setId")
		if policyId > 0 {
			pbPolicy, ok := wafPolicyCacheMap[policyId]
			if !ok {
				policyResp, err := this.RPC().HTTPFirewallPolicyRPC().FindEnabledHTTPFirewallPolicy(this.AdminContext(), &pb.FindEnabledHTTPFirewallPolicyRequest{HttpFirewallPolicyId: policyId})
				if err != nil {
					this.ErrorPage(err)
					return
				}
				pbPolicy = policyResp.HttpFirewallPolicy
				wafPolicyCacheMap[policyId] = pbPolicy
			}
			if pbPolicy != nil {
				wafMap = maps.Map{
					"policy": maps.Map{
						"id":       pbPolicy.Id,
						"name":     pbPolicy.Name,
						"serverId": pbPolicy.ServerId,
					},
				}
				if groupId > 0 {
					pbGroup, ok := wafGroupCacheMap[groupId]
					if !ok {
						groupResp, err := this.RPC().HTTPFirewallRuleGroupRPC().FindEnabledHTTPFirewallRuleGroup(this.AdminContext(), &pb.FindEnabledHTTPFirewallRuleGroupRequest{FirewallRuleGroupId: groupId})
						if err != nil {
							this.ErrorPage(err)
							return
						}
						pbGroup = groupResp.FirewallRuleGroup
						wafGroupCacheMap[groupId] = pbGroup
					}

					if pbGroup != nil {
						wafMap["group"] = maps.Map{
							"id":   pbGroup.Id,
							"name": pbGroup.Name,
						}

						if setId > 0 {
							pbSet, ok := wafSetCacheMap[setId]
							if !ok {
								setResp, err := this.RPC().HTTPFirewallRuleSetRPC().FindEnabledHTTPFirewallRuleSet(this.AdminContext(), &pb.FindEnabledHTTPFirewallRuleSetRequest{FirewallRuleSetId: setId})
								if err != nil {
									this.ErrorPage(err)
									return
								}
								pbSet = setResp.FirewallRuleSet
								wafSetCacheMap[setId] = pbSet
							}

							if pbSet != nil {
								wafMap["set"] = maps.Map{
									"id":   pbSet.Id,
									"name": pbSet.Name,
								}
							}
						}
					}
				}
			}
		}

		wafInfos[setId] = wafMap
	}
	this.Data["wafInfos"] = wafInfos

	this.Show()
}
