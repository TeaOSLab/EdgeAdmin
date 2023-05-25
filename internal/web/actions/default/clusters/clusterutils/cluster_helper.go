package clusterutils

import (
	teaconst "github.com/TeaOSLab/EdgeAdmin/internal/const"
	"github.com/TeaOSLab/EdgeAdmin/internal/utils/numberutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/dao"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"net/http"
	"strconv"
)

// ClusterHelper 单个集群的帮助
type ClusterHelper struct {
}

func NewClusterHelper() *ClusterHelper {
	return &ClusterHelper{}
}

func (this *ClusterHelper) BeforeAction(actionPtr actions.ActionWrapper) (goNext bool) {
	var action = actionPtr.Object()
	if action.Request.Method != http.MethodGet {
		return true
	}

	action.Data["teaMenu"] = "clusters"

	selectedTabbar := action.Data.GetString("mainTab")
	clusterId := action.ParamInt64("clusterId")
	clusterIdString := strconv.FormatInt(clusterId, 10)
	action.Data["clusterId"] = clusterId

	if clusterId > 0 {
		var ctx = actionPtr.(actionutils.ActionInterface).AdminContext()
		cluster, err := dao.SharedNodeClusterDAO.FindEnabledNodeCluster(ctx, clusterId)
		if err != nil {
			logs.Error(err)
			return
		}
		if cluster == nil {
			action.WriteString("can not find cluster")
			return
		}
		action.Data["currentClusterName"] = cluster.Name

		clusterInfo, err := dao.SharedNodeClusterDAO.FindEnabledNodeClusterConfigInfo(ctx, clusterId)
		if err != nil {
			logs.Error(err)
			return
		}
		if clusterInfo == nil {
			action.WriteString("can not find cluster info")
			return
		}

		var tabbar = actionutils.NewTabbar()
		tabbar.Add("", "", "/clusters", "arrow left", false)
		{
			var item = tabbar.Add(cluster.Name, "", "/clusters/cluster?clusterId="+clusterIdString, "angle right", true)
			item["isTitle"] = true
		}
		if teaconst.IsPlus {
			tabbar.Add("集群看板", "", "/clusters/cluster/boards?clusterId="+clusterIdString, "chart line area", selectedTabbar == "board")
		}
		tabbar.Add("集群节点", "", "/clusters/cluster/nodes?clusterId="+clusterIdString, "server", selectedTabbar == "node")
		tabbar.Add("集群设置", "", "/clusters/cluster/settings?clusterId="+clusterIdString, "setting", selectedTabbar == "setting")
		tabbar.Add("删除集群", "", "/clusters/cluster/delete?clusterId="+clusterIdString, "trash", selectedTabbar == "delete")
		actionutils.SetTabbar(action, tabbar)

		// 左侧菜单
		var secondMenuItem = action.Data.GetString("secondMenuItem")
		switch selectedTabbar {
		case "setting":
			var menuItems = this.createSettingMenu(cluster, clusterInfo, secondMenuItem)
			action.Data["leftMenuItems"] = menuItems

			// 当前菜单
			action.Data["leftMenuActiveItem"] = nil
			for _, item := range menuItems {
				if item.GetBool("isActive") {
					action.Data["leftMenuActiveItem"] = item
					break
				}
			}
		}
	}

	return true
}

// 设置菜单
func (this *ClusterHelper) createSettingMenu(cluster *pb.NodeCluster, info *pb.FindEnabledNodeClusterConfigInfoResponse, selectedItem string) (items []maps.Map) {
	clusterId := numberutils.FormatInt64(cluster.Id)
	items = append(items, maps.Map{
		"name":     "基础设置",
		"url":      "/clusters/cluster/settings?clusterId=" + clusterId,
		"isActive": selectedItem == "basic",
		"isOn":     true,
	})

	items = append(items, maps.Map{
		"name":     "DNS设置",
		"url":      "/clusters/cluster/settings/dns?clusterId=" + clusterId,
		"isActive": selectedItem == "dns",
		"isOn":     cluster.DnsDomainId > 0 || len(cluster.DnsName) > 0,
	})
	items = append(items, maps.Map{
		"name":     "健康检查",
		"url":      "/clusters/cluster/settings/health?clusterId=" + clusterId,
		"isActive": selectedItem == "health",
		"isOn":     info != nil && info.HealthCheckIsOn,
	})

	items = append(items, maps.Map{
		"name": "-",
	})

	items = append(items, maps.Map{
		"name":     "网站设置",
		"url":      "/clusters/cluster/settings/global-server-config?clusterId=" + clusterId,
		"isActive": selectedItem == "globalServerConfig",
		"isOn":     true,
	})

	items = append(items, maps.Map{
		"name":     "缓存策略",
		"url":      "/clusters/cluster/settings/cache?clusterId=" + clusterId,
		"isActive": selectedItem == "cache",
		"isOn":     cluster.HttpCachePolicyId > 0,
	})
	items = append(items, maps.Map{
		"name":     "WAF策略",
		"url":      "/clusters/cluster/settings/waf?clusterId=" + clusterId,
		"isActive": selectedItem == "waf",
		"isOn":     cluster.HttpFirewallPolicyId > 0,
	})

	items = append(items, maps.Map{
		"name":     "WAF动作",
		"url":      "/clusters/cluster/settings/firewall-actions?clusterId=" + clusterId,
		"isActive": selectedItem == "firewallAction",
		"isOn":     info != nil && info.HasFirewallActions,
	})

	items = append(items, maps.Map{
		"name":     "WebP",
		"url":      "/clusters/cluster/settings/webp?clusterId=" + clusterId,
		"isActive": selectedItem == "webp",
		"isOn":     info != nil && info.WebpIsOn,
	})

	items = filterMenuItems1(items, info, clusterId, selectedItem)

	items = append(items, maps.Map{
		"name":     "-",
		"url":      "",
		"isActive": false,
	})

	items = append(items, maps.Map{
		"name":     "DDoS防护",
		"url":      "/clusters/cluster/settings/ddos-protection?clusterId=" + clusterId,
		"isActive": selectedItem == "ddosProtection",
		"isOn":     info != nil && info.HasDDoSProtection,
	})

	items = append(items, maps.Map{
		"name": "-",
	})

	items = append(items, maps.Map{
		"name":     "统计指标",
		"url":      "/clusters/cluster/settings/metrics?clusterId=" + clusterId,
		"isActive": selectedItem == "metric",
		"isOn":     info != nil && info.HasMetricItems,
	})

	items = filterMenuItems2(items, info, clusterId, selectedItem)

	items = append(items, maps.Map{
		"name":     "-",
		"url":      "",
		"isActive": false,
	})

	items = append(items, maps.Map{
		"name":     "系统服务",
		"url":      "/clusters/cluster/settings/services?clusterId=" + clusterId,
		"isActive": selectedItem == "service",
		"isOn":     info != nil && info.HasSystemServices,
	})
	{
		items = append(items, maps.Map{
			"name":     "TOA设置",
			"url":      "/clusters/cluster/settings/toa?clusterId=" + clusterId,
			"isActive": selectedItem == "toa",
			"isOn":     info != nil && info.IsTOAEnabled,
		})
	}
	return
}
