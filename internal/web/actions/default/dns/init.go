package dns

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/dns/providers"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/helpers"
	"github.com/iwind/TeaGo"
)

func init() {
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		server.
			Helper(new(helpers.UserMustAuth)).
			Helper(new(Helper)).
			Prefix("/dns").
			Get("", new(IndexAction)).

			Prefix("/dns/providers").
			Data("teaSubMenu", "provider").
			Get("", new(providers.IndexAction)).
			GetPost("/createPopup", new(providers.CreatePopupAction)).
			GetPost("/updatePopup", new(providers.UpdatePopupAction)).
			Post("/delete", new(providers.DeleteAction)).
			EndData().

			EndAll()
	})
}
