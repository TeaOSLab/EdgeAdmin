package api

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/api/node"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/settings/settingutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/helpers"
	"github.com/iwind/TeaGo"
)

func init() {
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		server.
			Helper(helpers.NewUserMustAuth()).
			Helper(NewHelper()).
			Helper(settingutils.NewHelper("api")).
			Prefix("/api").
			Get("", new(IndexAction)).
			GetPost("/node/createPopup", new(node.CreatePopupAction)).
			Post("/delete", new(DeleteAction)).
			EndAll()
	})
}
