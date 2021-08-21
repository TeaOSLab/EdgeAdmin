package clusterutils

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/rpc"
	"github.com/TeaOSLab/EdgeAdmin/internal/utils/numberutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
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
	action := actionPtr.Object()
	if action.Request.Method != http.MethodGet {
		return true
	}

	action.Data["teaMenu"] = "ns"

	selectedTabbar := action.Data.GetString("mainTab")
	clusterId := action.ParamInt64("clusterId")
	clusterIdString := strconv.FormatInt(clusterId, 10)
	action.Data["clusterId"] = clusterId

	if clusterId > 0 {
		rpcClient, err := rpc.SharedRPC()
		if err != nil {
			logs.Error(err)
			return
		}
		clusterResp, err := rpcClient.NSClusterRPC().FindEnabledNSCluster(actionPtr.(actionutils.ActionInterface).AdminContext(), &pb.FindEnabledNSClusterRequest{
			NsClusterId: clusterId,
		})
		if err != nil {
			logs.Error(err)
			return
		}
		cluster := clusterResp.NsCluster
		if cluster == nil {
			action.WriteString("can not find ns cluster")
			return
		}

		tabbar := actionutils.NewTabbar()
		tabbar.Add("集群列表", "", "/ns/clusters", "", false)
		tabbar.Add("集群节点", "", "/ns/clusters/cluster?clusterId="+clusterIdString, "server", selectedTabbar == "node")
		tabbar.Add("集群设置", "", "/ns/clusters/cluster/settings?clusterId="+clusterIdString, "setting", selectedTabbar == "setting")
		tabbar.Add("删除集群", "", "/ns/clusters/cluster/delete?clusterId="+clusterIdString, "trash", selectedTabbar == "delete")

		{
			m := tabbar.Add("当前集群："+cluster.Name, "", "/ns/clusters/cluster?clusterId="+clusterIdString, "", false)
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

	return true
}

// 设置菜单
func (this *ClusterHelper) createSettingMenu(cluster *pb.NSCluster, selectedItem string) (items []maps.Map) {
	clusterId := numberutils.FormatInt64(cluster.Id)
	return []maps.Map{
		{
			"name":     "基础设置",
			"url":      "/ns/clusters/cluster/settings?clusterId=" + clusterId,
			"isActive": selectedItem == "basic",
		},
		{
			"name":     "访问日志",
			"url":      "/ns/clusters/cluster/settings/accessLog?clusterId=" + clusterId,
			"isActive": selectedItem == "accessLog",
		},
		{
			"name":     "递归DNS",
			"url":      "/ns/clusters/cluster/settings/recursion?clusterId=" + clusterId,
			"isActive": selectedItem == "recursion",
		},
	}
}
