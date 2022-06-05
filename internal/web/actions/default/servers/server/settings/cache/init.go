package cache

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/configloaders"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/servers/serverutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/helpers"
	"github.com/iwind/TeaGo"
)

func init() {
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		server.
			Helper(helpers.NewUserMustAuth(configloaders.AdminModuleCodeServer)).
			Helper(serverutils.NewServerHelper()).
			Prefix("/servers/server/settings/cache").
			GetPost("", new(IndexAction)).
			GetPost("/createPopup", new(CreatePopupAction)).
			GetPost("/purge", new(PurgeAction)).
			GetPost("/fetch", new(FetchAction)).
			Post("/updateRefs", new(UpdateRefsAction)).
			EndAll()
	})
}
