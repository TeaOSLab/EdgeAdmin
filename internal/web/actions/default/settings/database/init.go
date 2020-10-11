package profile

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/settings/settingutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/helpers"
	"github.com/iwind/TeaGo"
)

func init() {
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		server.
			Helper(helpers.NewUserMustAuth()).
			Helper(settingutils.NewHelper("database")).
			Prefix("/settings/database").
			Get("", new(IndexAction)).
			GetPost("/update", new(UpdateAction)).
			EndAll()
	})
}
