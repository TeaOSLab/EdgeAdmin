package reverseProxy

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
			Data("mainTab", "setting").
			Data("tinyMenuItem", "reverseProxy").
			Prefix("/servers/server/settings/locations/reverseProxy").
			GetPost("", new(IndexAction)).
			GetPost("/scheduling", new(SchedulingAction)).
			EndAll()
	})
}
