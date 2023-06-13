package updates

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
			Helper(settingutils.NewHelper("updates")).
			Prefix("/settings/updates").
			GetPost("", new(IndexAction)).
			Post("/update", new(UpdateAction)).
			Post("/ignoreVersion", new(IgnoreVersionAction)).
			Post("/resetIgnoredVersion", new(ResetIgnoredVersionAction)).
			GetPost("/upgrade", new(UpgradeAction)).
			EndAll()
	})
}
