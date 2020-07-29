package nodes

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/nodes/grants"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/helpers"
	"github.com/iwind/TeaGo"
)

func init() {
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		server.
			Helper(new(helpers.UserMustAuth)).
			Helper(new(Helper)).
			Prefix("/nodes").
			Get("", new(IndexAction)).
			GetPost("/create", new(CreateAction)).

			// 授权管理
			Get("/grants", new(grants.IndexAction)).
			GetPost("/grants/create", new(grants.CreateAction)).
			GetPost("/grants/update", new(grants.UpdateAction)).
			Post("/grants/delete", new(grants.DeleteAction)).
			Get("/grants/grant", new(grants.GrantAction)).
			EndAll()
	})
}
