package cluster

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/configloaders"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/clusters/cluster/groups"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/clusters/cluster/node"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/clusters/cluster/node/settings/cache"
	ddosProtection "github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/clusters/cluster/node/settings/ddos-protection"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/clusters/cluster/node/settings/dns"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/clusters/cluster/node/settings/ssh"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/clusters/cluster/node/settings/system"
	clusters "github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/clusters/clusterutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/helpers"
	"github.com/iwind/TeaGo"
)

func init() {
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		server.
			Helper(helpers.NewUserMustAuth(configloaders.AdminModuleCodeNode)).
			Helper(clusters.NewClusterHelper()).
			Data("teaMenu", "clusters").
			Data("teaSubMenu", "cluster").
			Prefix("/clusters/cluster").
			Get("", new(IndexAction)).
			Get("/nodes", new(NodesAction)).
			GetPost("/installNodes", new(InstallNodesAction)).
			GetPost("/installRemote", new(InstallRemoteAction)).
			Post("/installStatus", new(InstallStatusAction)).
			GetPost("/upgradeRemote", new(UpgradeRemoteAction)).
			Post("/upgradeStatus", new(UpgradeStatusAction)).
			GetPost("/delete", new(DeleteAction)).
			GetPost("/createNode", new(CreateNodeAction)).
			Post("/createNodeInstall", new(CreateNodeInstallAction)).
			GetPost("/createBatch", new(CreateBatchAction)).
			GetPost("/updateNodeSSH", new(UpdateNodeSSHAction)).
			GetPost("/installManual", new(InstallManualAction)).
			Post("/suggestLoginPorts", new(SuggestLoginPortsAction)).
			Get("/downloadInstaller", new(DownloadInstallerAction)).

			// 节点相关
			Prefix("/clusters/cluster/node").
			Get("", new(node.IndexAction)).
			GetPost("/update", new(node.UpdateAction)).
			GetPost("/install", new(node.InstallAction)).
			Post("/updateInstallStatus", new(node.UpdateInstallStatusAction)).
			Post("/status", new(node.StatusAction)).
			Get("/logs", new(node.LogsAction)).
			Post("/start", new(node.StartAction)).
			Post("/stop", new(node.StopAction)).
			Post("/uninstall", new(node.UninstallAction)).
			Post("/up", new(node.UpAction)).
			Post("/updateIsOn", new(node.UpdateIsOnAction)).
			Get("/detail", new(node.DetailAction)).
			GetPost("/updateDNSPopup", new(node.UpdateDNSPopupAction)).
			Post("/syncDomain", new(node.SyncDomainAction)).
			GetPost("/settings/cache", new(cache.IndexAction)).
			GetPost("/settings/dns", new(dns.IndexAction)).
			GetPost("/settings/system", new(system.IndexAction)).
			GetPost("/settings/ssh", new(ssh.IndexAction)).
			GetPost("/settings/ssh/test", new(ssh.TestAction)).
			GetPost("/settings/ddos-protection", new(ddosProtection.IndexAction)).
			Post("/settings/ddos-protection/status", new(ddosProtection.StatusAction)).

			// 分组相关
			Prefix("/clusters/cluster/groups").
			Get("", new(groups.IndexAction)).
			GetPost("/createPopup", new(groups.CreatePopupAction)).
			GetPost("/updatePopup", new(groups.UpdatePopupAction)).
			Post("/delete", new(groups.DeleteAction)).
			Post("/sort", new(groups.SortAction)).
			GetPost("/selectPopup", new(groups.SelectPopupAction)).
			EndAll()
	})
}
