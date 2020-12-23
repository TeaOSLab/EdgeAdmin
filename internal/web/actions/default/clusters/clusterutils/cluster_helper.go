package clusterutils

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/rpc"
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

// 单个集群的帮助
type ClusterHelper struct {
}

func NewClusterHelper() *ClusterHelper {
	return &ClusterHelper{}
}

func (this *ClusterHelper) BeforeAction(actionPtr actions.ActionWrapper) {
	action := actionPtr.Object()
	if action.Request.Method != http.MethodGet {
		return
	}

	action.Data["teaMenu"] = "clusters"

	selectedTabbar := action.Data.GetString("mainTab")
	clusterId := action.ParamInt64("clusterId")
	clusterIdString := strconv.FormatInt(clusterId, 10)
	action.Data["clusterId"] = clusterId

	if clusterId > 0 {
		cluster, err := dao.SharedNodeClusterDAO.FindEnabledNodeCluster(actionPtr.(rpc.ContextInterface).AdminContext(), clusterId)
		if err != nil {
			logs.Error(err)
			return
		}
		if cluster == nil {
			action.WriteString("can not find cluster")
			return
		}

		tabbar := actionutils.NewTabbar()
		tabbar.Add("集群列表", "", "/clusters", "", false)
		tabbar.Add("集群节点", "", "/clusters/cluster?clusterId="+clusterIdString, "server", selectedTabbar == "node")
		tabbar.Add("集群设置", "", "/clusters/cluster/settings?clusterId="+clusterIdString, "setting", selectedTabbar == "setting")
		tabbar.Add("删除集群", "", "/clusters/cluster/delete?clusterId="+clusterIdString, "trash", selectedTabbar == "delete")

		{
			m := tabbar.Add("当前集群："+cluster.Name, "", "/clusters/cluster?clusterId="+clusterIdString, "", false)
			m["right"] = true
		}
		actionutils.SetTabbar(action, tabbar)

		// 左侧菜单
		secondMenuItem := action.Data.GetString("secondMenuItem")
		switch selectedTabbar {
		case "setting":
			action.Data["leftMenuItems"] = this.createSettingMenu(cluster, secondMenuItem)
		}
	}
}

// 设置菜单
func (this *ClusterHelper) createSettingMenu(cluster *pb.NodeCluster, selectedItem string) (items []maps.Map) {
	clusterId := numberutils.FormatInt64(cluster.Id)
	items = append(items, maps.Map{
		"name":     "基础设置",
		"url":      "/clusters/cluster/settings?clusterId=" + clusterId,
		"isActive": selectedItem == "basic",
	})
	items = append(items, maps.Map{
		"name":     "缓存设置",
		"url":      "/clusters/cluster/settings/cache?clusterId=" + clusterId,
		"isActive": selectedItem == "cache",
		"isOn":     cluster.HttpCachePolicyId > 0,
	})
	items = append(items, maps.Map{
		"name":     "WAF设置",
		"url":      "/clusters/cluster/settings/waf?clusterId=" + clusterId,
		"isActive": selectedItem == "waf",
		"isOn":     cluster.HttpFirewallPolicyId > 0,
	})
	items = append(items, maps.Map{
		"name":     "健康检查",
		"url":      "/clusters/cluster/settings/health?clusterId=" + clusterId,
		"isActive": selectedItem == "health",
	})
	items = append(items, maps.Map{
		"name":     "DNS设置",
		"url":      "/clusters/cluster/settings/dns?clusterId=" + clusterId,
		"isActive": selectedItem == "dns",
	})
	items = append(items, maps.Map{
		"name":     "TOA设置",
		"url":      "/clusters/cluster/settings/toa?clusterId=" + clusterId,
		"isActive": selectedItem == "toa",
	})
	return
}
