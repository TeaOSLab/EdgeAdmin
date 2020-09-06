package nodes

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/nodes/grants"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/nodes/ipAddresses"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/helpers"
	"github.com/iwind/TeaGo"
)

func init() {
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		server.
			Helper(new(helpers.UserMustAuth)).
			Helper(new(Helper)).
			Prefix("/nodes").
			Post("/delete", new(DeleteAction)).

			// IP地址
			GetPost("/ipAddresses/createPopup", new(ipAddresses.CreatePopupAction)).
			GetPost("/ipAddresses/updatePopup", new(ipAddresses.UpdatePopupAction)).

			// 授权管理
			Get("/grants", new(grants.IndexAction)).
			GetPost("/grants/create", new(grants.CreateAction)).
			GetPost("/grants/update", new(grants.UpdateAction)).
			Post("/grants/delete", new(grants.DeleteAction)).
			Get("/grants/grant", new(grants.GrantAction)).
			GetPost("/grants/selectPopup", new(grants.SelectPopupAction)).
			GetPost("/grants/createPopup", new(grants.CreatePopupAction)).
			GetPost("/grants/updatePopup", new(grants.UpdatePopupAction)).
			EndAll()
	})
}
