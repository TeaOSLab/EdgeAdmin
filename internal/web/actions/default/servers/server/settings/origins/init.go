package origins

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/configloaders"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/helpers"
	"github.com/iwind/TeaGo"
)

func init() {
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		server.
			Helper(helpers.NewUserMustAuth(configloaders.AdminModuleCodeServer)).
			Prefix("/servers/server/settings/origins").
			GetPost("/addPopup", new(AddPopupAction)).
			Post("/delete", new(DeleteAction)).
			GetPost("/updatePopup", new(UpdatePopupAction)).
			Post("/updateIsOn", new(UpdateIsOnAction)).
			Post("/detectHTTPS", new(DetectHTTPSAction)).
			EndAll()
	})
}
