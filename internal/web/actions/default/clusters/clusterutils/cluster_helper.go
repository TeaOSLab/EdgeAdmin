package clusters

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/rpc"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"net/http"
	"strconv"
)

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
	tabbar.Add("当前集群："+cluster.Name, "", "/clusters", "left long alternate arrow", false)
	tabbar.Add("节点", "", "/clusters/cluster?clusterId="+clusterIdString, "server", selectedTabbar == "node")
	tabbar.Add("设置", "", "/clusters/cluster/settings?clusterId="+clusterIdString, "setting", selectedTabbar == "setting")
	actionutils.SetTabbar(action, tabbar)

	// 左侧菜单
	secondMenuItem := action.Data.GetString("secondMenuItem")
	switch selectedTabbar {
	case "setting":
		action.Data["leftMenuItems"] = this.createSettingMenu(clusterIdString, secondMenuItem)
	case "node":
		action.Data["leftMenuItems"] = this.createNodeMenu(clusterIdString, secondMenuItem)
	}
}

// 节点菜单
func (this *ClusterHelper) createNodeMenu(clusterId string, selectedItem string) (items []maps.Map) {
	items = append(items, maps.Map{
		"name":     "节点列表",
		"url":      "/clusters/cluster?clusterId=" + clusterId,
		"isActive": selectedItem == "nodes",
	})
	return
}

// 设置菜单
func (this *ClusterHelper) createSettingMenu(clusterId string, selectedItem string) (items []maps.Map) {
	items = append(items, maps.Map{
		"name":     "基础设置",
		"url":      "/clusters/cluster/settings?clusterId=" + clusterId,
		"isActive": selectedItem == "basic",
	})
	return
}
