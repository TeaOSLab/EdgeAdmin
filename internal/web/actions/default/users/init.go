package users

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/configloaders"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/users/accesskeys"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/helpers"
	"github.com/iwind/TeaGo"
)

func init() {
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		server.
			Helper(helpers.NewUserMustAuth(configloaders.AdminModuleCodeUser)).
			Data("teaMenu", "users").
			Prefix("/users").
			Data("teaSubMenu", "users").
			Get("", new(IndexAction)).
			GetPost("/createPopup", new(CreatePopupAction)).

			// 单个用户信息
			Get("/user", new(UserAction)).
			GetPost("/update", new(UpdateAction)).
			Post("/delete", new(DeleteAction)).
			GetPost("/features", new(FeaturesAction)).
			GetPost("/verifyPopup", new(VerifyPopupAction)).
			Get("/otpQrcode", new(OtpQrcodeAction)).

			// AccessKeys
			Prefix("/users/accesskeys").
			Get("", new(accesskeys.IndexAction)).
			GetPost("/createPopup", new(accesskeys.CreatePopupAction)).
			Post("/delete", new(accesskeys.DeleteAction)).
			Post("/updateIsOn", new(accesskeys.UpdateIsOnAction)).

			//
			EndAll()
	})
}
