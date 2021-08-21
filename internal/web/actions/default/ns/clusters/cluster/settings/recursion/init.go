package recursion

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/configloaders"
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
			Prefix("/ns/clusters/cluster/settings/recursion").
			GetPost("", new(IndexAction)).
			EndAll()
	})
}
