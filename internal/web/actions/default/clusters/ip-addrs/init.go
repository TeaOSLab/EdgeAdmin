package ipaddrs

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/configloaders"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/clusters/clusterutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/clusters/ip-addrs/addr"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/helpers"
	"github.com/iwind/TeaGo"
)

func init() {
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		server.
			Helper(helpers.NewUserMustAuth(configloaders.AdminModuleCodeNode)).
			Helper(clusterutils.NewClustersHelper()).
			Data("teaSubMenu", "ipAddr").
			Prefix("/clusters/ip-addrs").
			Get("", new(IndexAction)).
			Get("/logs", new(LogsAction)).

			// 单个地址操作
			Post("/addr/delete", new(addr.DeleteAction)).
			Get("/addr", new(addr.IndexAction)).
			GetPost("/addr/update", new(addr.UpdateAction)).
			Get("/addr/logs", new(addr.LogsAction)).
			EndAll()
	})
}
