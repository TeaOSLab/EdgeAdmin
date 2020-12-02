package admins

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/helpers"
	"github.com/iwind/TeaGo"
)

func init() {
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		server.
			Helper(helpers.NewUserMustAuth()).
			Data("teaMenu", "admins").
			Prefix("/admins").
			Get("", new(IndexAction)).
			GetPost("/createPopup", new(CreatePopupAction)).
			GetPost("/updatePopup", new(UpdatePopupAction)).
			Post("/delete", new(DeleteAction)).
			Post("/updateOn", new(UpdateOnAction)).
			EndAll()
	})
}
