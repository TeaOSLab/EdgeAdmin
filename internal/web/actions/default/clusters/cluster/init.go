package cluster

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/clusters/cluster/groups"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/clusters/cluster/node"
	clusters "github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/clusters/clusterutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/helpers"
	"github.com/iwind/TeaGo"
)

func init() {
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		server.
			Helper(helpers.NewUserMustAuth()).
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

			// 分组相关
			Get("/groups", new(groups.IndexAction)).
			GetPost("/groups/createPopup", new(groups.CreatePopupAction)).
			GetPost("/groups/updatePopup", new(groups.UpdatePopupAction)).
			Post("/groups/delete", new(groups.DeleteAction)).
			Post("/groups/sort", new(groups.SortAction)).
			Post("/groups/options", new(groups.OptionsAction)).
			GetPost("/groups/selectPopup", new(groups.SelectPopupAction)).

			EndAll()
	})
}
