package servers

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/helpers"
	"github.com/iwind/TeaGo"
)

func init() {
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		server.
			Helper(helpers.NewUserMustAuth()).
			Helper(NewHelper()).
			Data("teaModule", "server").
			Prefix("/servers").
			Get("", new(IndexAction)).
			GetPost("/create", new(CreateAction)).
			GetPost("/update", new(UpdateAction)).

			GetPost("/addPortPopup", new(AddPortPopupAction)).
			GetPost("/addServerNamePopup", new(AddServerNamePopupAction)).
			GetPost("/addOriginPopup", new(AddOriginPopupAction)).
			EndAll()
	})
}
