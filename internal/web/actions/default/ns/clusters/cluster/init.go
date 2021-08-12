package cluster

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/configloaders"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/ns/clusters/cluster/node"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/ns/clusters/clusterutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/helpers"
	"github.com/iwind/TeaGo"
)

func init() {
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		server.
			Helper(helpers.NewUserMustAuth(configloaders.AdminModuleCodeNS)).
			Helper(new(clusterutils.ClusterHelper)).
			Data("teaMenu", "ns").
			Data("teaSubMenu", "cluster").
			Prefix("/ns/clusters/cluster").
			Get("", new(IndexAction)).
			GetPost("/delete", new(DeleteAction)).
			GetPost("/createNode", new(CreateNodeAction)).
			Post("/deleteNode", new(DeleteNodeAction)).
			Get("/upgradeRemote", new(UpgradeRemoteAction)).
			GetPost("/updateNodeSSH", new(UpdateNodeSSHAction)).

			// 节点相关
			Prefix("/ns/clusters/cluster/node").
			Get("", new(node.IndexAction)).
			Get("/logs", new(node.LogsAction)).
			GetPost("/update", new(node.UpdateAction)).
			GetPost("/install", new(node.InstallAction)).
			Post("/status", new(node.StatusAction)).
			Post("/updateInstallStatus", new(node.UpdateInstallStatusAction)).
			Post("/start", new(node.StartAction)).
			Post("/stop", new(node.StopAction)).
			EndAll()
	})
}
