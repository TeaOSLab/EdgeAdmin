package ns

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/configloaders"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/ns/domains"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/ns/domains/keys"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/ns/domains/records"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/helpers"
	"github.com/iwind/TeaGo"
)

func init() {
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		server.
			Helper(helpers.NewUserMustAuth(configloaders.AdminModuleCodeNS)).
			Data("teaMenu", "ns").
			Prefix("/ns").
			Get("", new(IndexAction)).

			// 域名相关
			Prefix("/ns/domains").
			GetPost("/create", new(domains.CreateAction)).
			Post("/delete", new(domains.DeleteAction)).
			Get("/domain", new(domains.DomainAction)).
			GetPost("/update", new(domains.UpdateAction)).

			// 域名密钥
			Prefix("/ns/domains/keys").
			Get("", new(keys.IndexAction)).
			GetPost("/createPopup", new(keys.CreatePopupAction)).
			GetPost("/updatePopup", new(keys.UpdatePopupAction)).
			Post("/delete", new(keys.DeleteAction)).
			Post("/generateSecret", new(keys.GenerateSecretAction)).

			// 记录相关
			Prefix("/ns/domains/records").
			Get("", new(records.IndexAction)).
			GetPost("/createPopup", new(records.CreatePopupAction)).
			GetPost("/updatePopup", new(records.UpdatePopupAction)).
			Post("/delete", new(records.DeleteAction)).

			EndAll()
	})
}
