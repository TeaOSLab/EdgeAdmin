package node

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/configloaders"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/settings/settingutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/helpers"
	"github.com/iwind/TeaGo"
)

func init() {
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		server.
			Helper(helpers.NewUserMustAuth(configloaders.AdminModuleCodeSetting)).
			Helper(settingutils.NewAdvancedHelper("monitorNodes")).
			Prefix("/settings/monitorNodes/node").

			// 节点相关
			Helper(NewHelper()).
			Get("", new(IndexAction)).
			Get("/logs", new(LogsAction)).
			GetPost("/update", new(UpdateAction)).
			Get("/install", new(InstallAction)).

			EndAll()
	})
}
