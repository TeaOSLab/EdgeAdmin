package node

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/settings/settingutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/helpers"
	"github.com/iwind/TeaGo"
)

func init() {
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		server.
			Helper(helpers.NewUserMustAuth()).
			Helper(settingutils.NewHelper("api")).
			Prefix("/api/node").

			// 这里不受Helper的约束
			GetPost("/createAddrPopup", new(CreateAddrPopupAction)).
			GetPost("/updateAddrPopup", new(UpdateAddrPopupAction)).

			// 节点相关
			Helper(NewHelper()).
			GetPost("/settings", new(SettingsAction)).


			EndAll()
	})
}
