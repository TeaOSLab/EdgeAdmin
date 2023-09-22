package log

import (
	"github.com/TeaOSLab/EdgeCommon/pkg/iplibrary"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/dao"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/maps"
	timeutil "github.com/iwind/TeaGo/utils/time"
)

type IndexAction struct {
	BaseAction
}

func (this *IndexAction) Init() {
	this.Nav("", "log", "")
	this.SecondMenu("index")
}

func (this *IndexAction) RunGet(params struct {
	ServerId  int64
	RequestId string
	Ip        string
	Domain    string
	ClusterId int64
	NodeId    int64
	Keyword   string
}) {
	this.Data["serverId"] = params.ServerId
	this.Data["requestId"] = params.RequestId
	this.Data["ip"] = params.Ip
	this.Data["domain"] = params.Domain
	this.Data["keyword"] = params.Keyword
	this.Data["path"] = this.Request.URL.Path
	this.Data["clusterId"] = params.ClusterId
	this.Data["nodeId"] = params.NodeId

	// 检查集群全局设置
	if !this.initClusterAccessLogConfig(params.ServerId) {
		return
	}

	// 检查当前网站有无开启访问日志
	this.Data["serverAccessLogIsOn"] = true

	groupResp, err := this.RPC().ServerGroupRPC().FindEnabledServerGroupConfigInfo(this.AdminContext(), &pb.FindEnabledServerGroupConfigInfoRequest{
		ServerId: params.ServerId,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	if !groupResp.HasAccessLogConfig {
		webConfig, err := dao.SharedHTTPWebDAO.FindWebConfigWithServerId(this.AdminContext(), params.ServerId)
		if err != nil {
			this.ErrorPage(err)
			return
		}
		if webConfig != nil && webConfig.AccessLogRef != nil && !webConfig.AccessLogRef.IsOn {
			this.Data["serverAccessLogIsOn"] = false
		}
	}

	// 记录最近使用
	_, err = this.RPC().LatestItemRPC().IncreaseLatestItem(this.AdminContext(), &pb.IncreaseLatestItemRequest{
		ItemType: "server",
		ItemId:   params.ServerId,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Show()
}

func (this *IndexAction) RunPost(params struct {
	ServerId  int64
	RequestId string
	Keyword   string
	Ip        string
	Domain    string
	ClusterId int64
	NodeId    int64

	Partition int32 `default:"-1"`

	Must *actions.Must
}) {
	var isReverse = len(params.RequestId) > 0
	accessLogsResp, err := this.RPC().HTTPAccessLogRPC().ListHTTPAccessLogs(this.AdminContext(), &pb.ListHTTPAccessLogsRequest{
		Partition:     params.Partition,
		ServerId:      params.ServerId,
		RequestId:     params.RequestId,
		Size:          20,
		Day:           timeutil.Format("Ymd"),
		Keyword:       params.Keyword,
		Ip:            params.Ip,
		Domain:        params.Domain,
		NodeId:        params.NodeId,
		NodeClusterId: params.ClusterId,
		Reverse:       isReverse,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	var ipList = []string{}
	var wafMaps = []maps.Map{}

	var accessLogs = accessLogsResp.HttpAccessLogs
	if len(accessLogs) == 0 {
		accessLogs = []*pb.HTTPAccessLog{}
	} else {
		for _, accessLog := range accessLogs {
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
	this.Data["accessLogs"] = accessLogs
	if len(accessLogs) > 0 {
		this.Data["requestId"] = accessLogs[0].RequestId
	} else {
		this.Data["requestId"] = params.RequestId
	}
	this.Data["hasMore"] = accessLogsResp.HasMore

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

	this.Success()
}
