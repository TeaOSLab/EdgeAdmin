package clusterutils

import (
	teaconst "github.com/TeaOSLab/EdgeAdmin/internal/const"
	"github.com/TeaOSLab/EdgeAdmin/internal/utils/numberutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/helpers"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
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
	helpers.LangHelper
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

	var selectedTabbar = action.Data.GetString("mainTab")
	var clusterId = action.ParamInt64("clusterId")
	var clusterIdString = strconv.FormatInt(clusterId, 10)
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

		var nodeId = action.ParamInt64("nodeId")
		var isInCluster = nodeId <= 0

		var tabbar = actionutils.NewTabbar()
		{
			var url = "/clusters"
			if !isInCluster {
				url = "/clusters/cluster/nodes?clusterId=" + clusterIdString
			}
			tabbar.Add("", "", url, "arrow left", false)
		}
		{
			var url = "/clusters/cluster?clusterId=" + clusterIdString
			if !isInCluster {
				url = "/clusters/cluster/nodes?clusterId=" + clusterIdString
			}

			var item = tabbar.Add(cluster.Name, "", url, "angle right", true)
			item.IsTitle = true
		}
		if teaconst.IsPlus {
			{
				var item = tabbar.Add(this.Lang(actionPtr, codes.NodeClusterMenu_TabClusterDashboard), "", "/clusters/cluster/boards?clusterId="+clusterIdString, "chart line area", selectedTabbar == "board")
				item.IsDisabled = !isInCluster
			}
		}
		{
			var item = tabbar.Add(this.Lang(actionPtr, codes.NodeClusterMenu_TabClusterNodes), "", "/clusters/cluster/nodes?clusterId="+clusterIdString, "server", selectedTabbar == "node")
			item.IsDisabled = !isInCluster
		}

		{
			var item = tabbar.Add(this.Lang(actionPtr, codes.NodeClusterMenu_TabClusterSettings), "", "/clusters/cluster/settings?clusterId="+clusterIdString, "setting", selectedTabbar == "setting")
			item.IsDisabled = !isInCluster
		}
		{
			var item = tabbar.Add(this.Lang(actionPtr, codes.NodeClusterMenu_TabClusterDelete), "", "/clusters/cluster/delete?clusterId="+clusterIdString, "trash", selectedTabbar == "delete")
			item.IsDisabled = !isInCluster
		}
		actionutils.SetTabbar(action, tabbar)

		// 左侧菜单
		var secondMenuItem = action.Data.GetString("secondMenuItem")
		switch selectedTabbar {
		case "setting":
			var menuItems = this.createSettingMenu(cluster, clusterInfo, secondMenuItem, actionPtr)
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
func (this *ClusterHelper) createSettingMenu(cluster *pb.NodeCluster, info *pb.FindEnabledNodeClusterConfigInfoResponse, selectedItem string, actionPtr actions.ActionWrapper) (items []maps.Map) {
	clusterId := numberutils.FormatInt64(cluster.Id)
	items = append(items, maps.Map{
		"name":     this.Lang(actionPtr, codes.NodeClusterMenu_SettingBasic),
		"url":      "/clusters/cluster/settings?clusterId=" + clusterId,
		"isActive": selectedItem == "basic",
		"isOn":     true,
	})

	items = append(items, maps.Map{
		"name":     this.Lang(actionPtr, codes.NodeClusterMenu_SettingDNS),
		"url":      "/clusters/cluster/settings/dns?clusterId=" + clusterId,
		"isActive": selectedItem == "dns",
		"isOn":     cluster.DnsDomainId > 0 || len(cluster.DnsName) > 0,
	})
	items = append(items, maps.Map{
		"name":     this.Lang(actionPtr, codes.NodeClusterMenu_SettingHealthCheck),
		"url":      "/clusters/cluster/settings/health?clusterId=" + clusterId,
		"isActive": selectedItem == "health",
		"isOn":     info != nil && info.HealthCheckIsOn,
	})

	items = append(items, maps.Map{
		"name": "-",
	})

	items = append(items, maps.Map{
		"name":     this.Lang(actionPtr, codes.NodeClusterMenu_SettingServiceGlobal),
		"url":      "/clusters/cluster/settings/global-server-config?clusterId=" + clusterId,
		"isActive": selectedItem == "globalServerConfig",
		"isOn":     true,
	})

	items = append(items, maps.Map{
		"name":     this.Lang(actionPtr, codes.NodeClusterMenu_SettingCachePolicy),
		"url":      "/clusters/cluster/settings/cache?clusterId=" + clusterId,
		"isActive": selectedItem == "cache",
		"isOn":     cluster.HttpCachePolicyId > 0,
	})
	items = append(items, maps.Map{
		"name":     this.Lang(actionPtr, codes.NodeClusterMenu_SettingWAFPolicy),
		"url":      "/clusters/cluster/settings/waf?clusterId=" + clusterId,
		"isActive": selectedItem == "waf",
		"isOn":     cluster.HttpFirewallPolicyId > 0,
	})

	items = append(items, maps.Map{
		"name":     this.Lang(actionPtr, codes.NodeClusterMenu_SettingWAFActions),
		"url":      "/clusters/cluster/settings/firewall-actions?clusterId=" + clusterId,
		"isActive": selectedItem == "firewallAction",
		"isOn":     info != nil && info.HasFirewallActions,
	})

	items = append(items, maps.Map{
		"name":     this.Lang(actionPtr, codes.NodeClusterMenu_SettingWebP),
		"url":      "/clusters/cluster/settings/webp?clusterId=" + clusterId,
		"isActive": selectedItem == "webp",
		"isOn":     info != nil && info.WebPIsOn,
	})

	items = this.filterMenuItems1(items, info, clusterId, selectedItem, actionPtr)

	items = append(items, maps.Map{
		"name":     this.Lang(actionPtr, codes.NodeClusterMenu_SettingMetrics),
		"url":      "/clusters/cluster/settings/metrics?clusterId=" + clusterId,
		"isActive": selectedItem == "metric",
		"isOn":     info != nil && info.HasMetricItems,
	})

	items = append(items, maps.Map{
		"name":     "-",
		"url":      "",
		"isActive": false,
	})

	items = append(items, maps.Map{
		"name":     this.Lang(actionPtr, codes.NodeClusterMenu_SettingDDoSProtection),
		"url":      "/clusters/cluster/settings/ddos-protection?clusterId=" + clusterId,
		"isActive": selectedItem == "ddosProtection",
		"isOn":     info != nil && info.HasDDoSProtection,
	})

	items = this.filterMenuItems2(items, info, clusterId, selectedItem, actionPtr)

	return
}
