package clusters

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/configloaders"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/clusters/clusterutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/helpers"
	"github.com/iwind/TeaGo"
)

func init() {
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		server.
			Helper(helpers.NewUserMustAuth(configloaders.AdminModuleCodeNode)).
			Helper(clusterutils.NewClustersHelper()).
			Data("teaMenu", "clusters").
			Data("teaSubMenu", "cluster").
			Prefix("/clusters").
			Get("", new(IndexAction)).
			GetPost("/create", new(CreateAction)).
			GetPost("/createNode", new(CreateNodeAction)).
			Post("/pin", new(PinAction)).
			Get("/nodes", new(NodesAction)).

			// 只要登录即可访问的Action
			EndHelpers().
			Helper(helpers.NewUserMustAuth(configloaders.AdminModuleCodeCommon)).
			Post("/options", new(OptionsAction)).
			Post("/nodeOptions", new(NodeOptionsAction)).
			GetPost("/selectPopup", new(SelectPopupAction)).
			EndAll()
	})
}
