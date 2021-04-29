package cluster

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/configloaders"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/clusters/cluster/groups"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/clusters/cluster/node"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/clusters/cluster/node/monitor"
	clusters "github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/clusters/clusterutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/helpers"
	"github.com/iwind/TeaGo"
)

func init() {
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		server.
			Helper(helpers.NewUserMustAuth(configloaders.AdminModuleCodeNode)).
			Helper(clusters.NewClusterHelper()).
			Prefix("/clusters/cluster").
			Get("", new(IndexAction)).
			GetPost("/installNodes", new(InstallNodesAction)).
			GetPost("/installRemote", new(InstallRemoteAction)).
			Post("/installStatus", new(InstallStatusAction)).
			GetPost("/upgradeRemote", new(UpgradeRemoteAction)).
			Post("/upgradeStatus", new(UpgradeStatusAction)).
			GetPost("/delete", new(DeleteAction)).
			GetPost("/createNode", new(CreateNodeAction)).
			GetPost("/createBatch", new(CreateBatchAction)).
			GetPost("/updateNodeSSH", new(UpdateNodeSSHAction)).
			GetPost("/installManual", new(InstallManualAction)).

			// 节点相关
			Get("/node", new(node.NodeAction)).
			GetPost("/node/update", new(node.UpdateAction)).
			GetPost("/node/install", new(node.InstallAction)).
			Post("/node/updateInstallStatus", new(node.UpdateInstallStatusAction)).
			Post("/node/status", new(node.StatusAction)).
			Get("/node/logs", new(node.LogsAction)).
			Post("/node/start", new(node.StartAction)).
			Post("/node/stop", new(node.StopAction)).
			Post("/node/up", new(node.UpAction)).
			Get("/node/monitor", new(monitor.IndexAction)).
			Post("/node/monitor/cpu", new(monitor.CpuAction)).
			Post("/node/monitor/memory", new(monitor.MemoryAction)).
			Post("/node/monitor/load", new(monitor.LoadAction)).
			Post("/node/monitor/trafficIn", new(monitor.TrafficInAction)).
			Post("/node/monitor/trafficOut", new(monitor.TrafficOutAction)).

			// 分组相关
			Get("/groups", new(groups.IndexAction)).
			GetPost("/groups/createPopup", new(groups.CreatePopupAction)).
			GetPost("/groups/updatePopup", new(groups.UpdatePopupAction)).
			Post("/groups/delete", new(groups.DeleteAction)).
			Post("/groups/sort", new(groups.SortAction)).
			GetPost("/groups/selectPopup", new(groups.SelectPopupAction)).

			EndAll()
	})
}
