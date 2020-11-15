package clusterutils

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/rpc"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
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

func (this *ClusterHelper) BeforeAction(action *actions.ActionObject) {
	if action.Request.Method != http.MethodGet {
		return
	}

	action.Data["teaMenu"] = "clusters"

	selectedTabbar := action.Data.GetString("mainTab")
	clusterId := action.ParamInt64("clusterId")
	clusterIdString := strconv.FormatInt(clusterId, 10)
	action.Data["clusterId"] = clusterId

	rpcClient, err := rpc.SharedRPC()
	if err != nil {
		logs.Error(err)
		return
	}

	if clusterId > 0 {
		clusterResp, err := rpcClient.NodeClusterRPC().FindEnabledNodeCluster(rpcClient.Context(action.Context.GetInt64("adminId")), &pb.FindEnabledNodeClusterRequest{ClusterId: clusterId})
		if err != nil {
			logs.Error(err)
			return
		}
		cluster := clusterResp.Cluster
		if cluster == nil {
			action.WriteString("can not find cluster")
			return
		}

		tabbar := actionutils.NewTabbar()
		tabbar.Add("集群列表", "", "/clusters", "", false)
		tabbar.Add("节点", "", "/clusters/cluster?clusterId="+clusterIdString, "server", selectedTabbar == "node")
		tabbar.Add("设置", "", "/clusters/cluster/settings?clusterId="+clusterIdString, "setting", selectedTabbar == "setting")
		tabbar.Add("删除", "", "/clusters/cluster/delete?clusterId="+clusterIdString, "trash", selectedTabbar == "delete")

		{
			m := tabbar.Add("当前集群："+cluster.Name, "", "/clusters/cluster?clusterId="+clusterIdString, "", false)
			m["right"] = true
		}
		actionutils.SetTabbar(action, tabbar)

		// 左侧菜单
		secondMenuItem := action.Data.GetString("secondMenuItem")
		switch selectedTabbar {
		case "setting":
			action.Data["leftMenuItems"] = this.createSettingMenu(clusterIdString, secondMenuItem)
		}
	}
}

// 设置菜单
func (this *ClusterHelper) createSettingMenu(clusterId string, selectedItem string) (items []maps.Map) {
	items = append(items, maps.Map{
		"name":     "基础设置",
		"url":      "/clusters/cluster/settings?clusterId=" + clusterId,
		"isActive": selectedItem == "basic",
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
	return
}
