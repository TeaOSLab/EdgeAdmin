package servers

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/configloaders"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/servers/headers"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/servers/users"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/helpers"
	"github.com/iwind/TeaGo"
)

func init() {
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		server.
			Data("teaMenu", "servers").
			Helper(helpers.NewUserMustAuth(configloaders.AdminModuleCodeServer)).
			Helper(NewHelper()).
			Prefix("/servers").
			Get("", new(IndexAction)).
			GetPost("/create", new(CreateAction)).
			GetPost("/update", new(UpdateAction)).
			Post("/nearby", new(NearbyAction)).
			Post("/deleteServers", new(DeleteServersAction)).

			//
			GetPost("/addPortPopup", new(AddPortPopupAction)).
			GetPost("/addServerNamePopup", new(AddServerNamePopupAction)).
			GetPost("/addOriginPopup", new(AddOriginPopupAction)).
			Get("/serverNamesPopup", new(ServerNamesPopupAction)).
			Post("/status", new(StatusAction)).

			// user
			Post("/users/options", new(users.OptionsAction)).
			Post("/users/plans", new(users.PlansAction)).

			// header
			Post("/headers/options", new(headers.OptionsAction)).

			//
			EndAll()
	})
}
