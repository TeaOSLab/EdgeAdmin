package cluster

import (
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

			// 节点相关
			Get("/node", new(node.NodeAction)).
			GetPost("/node/create", new(node.CreateAction)).
			GetPost("/node/update", new(node.UpdateAction)).
			GetPost("/node/install", new(node.InstallAction)).
			Post("/node/updateInstallStatus", new(node.UpdateInstallStatusAction)).
			Post("/node/status", new(node.StatusAction)).
			EndAll()
	})
}
