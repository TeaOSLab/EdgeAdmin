package reverseProxy

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
			Data("mainTab", "setting").
			Data("secondMenuItem", "reverseProxy").
			Prefix("/servers/server/settings/reverseProxy").
			Get("", new(IndexAction)).
			GetPost("/scheduling", new(SchedulingAction)).
			GetPost("/updateSchedulingPopup", new(UpdateSchedulingPopupAction)).
			GetPost("/setting", new(SettingAction)).
			EndAll()
	})
}
