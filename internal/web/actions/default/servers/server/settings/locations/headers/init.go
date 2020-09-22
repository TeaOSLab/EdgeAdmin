package headers

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/servers/server/settings/locations/locationutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/servers/serverutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/helpers"
	"github.com/iwind/TeaGo"
)

func init() {
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		server.
			Helper(helpers.NewUserMustAuth()).
			Helper(locationutils.NewLocationHelper()).
			Helper(serverutils.NewServerHelper()).
			Data("tinyMenuItem", "header").
			Prefix("/servers/server/settings/locations/headers").
			Get("", new(IndexAction)).
			EndAll()
	})
}
