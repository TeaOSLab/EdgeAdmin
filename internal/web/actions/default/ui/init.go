package ui

import (
	"compress/gzip"
	"github.com/TeaOSLab/EdgeAdmin/internal/configloaders"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/helpers"
	"github.com/iwind/TeaGo"
	"github.com/iwind/TeaGo/actions"
)

func init() {
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		server.
			Prefix("/ui").

			// 公共可以访问的链接
			Get("/image/:fileId", new(ImageAction)).

			// 以下的需要压缩
			Helper(&actions.Gzip{Level: gzip.BestCompression}).
			Get("/components.js", new(ComponentsAction)).
			EndHelpers().

			// 以下需要登录
			Helper(helpers.NewUserMustAuth(configloaders.AdminModuleCodeCommon)).
			Get("/download", new(DownloadAction)).
			GetPost("/selectProvincesPopup", new(SelectProvincesPopupAction)).
			GetPost("/selectCountriesPopup", new(SelectCountriesPopupAction)).
			Post("/eventLevelOptions", new(EventLevelOptionsAction)).
			Post("/showTip", new(ShowTipAction)).
			Post("/hideTip", new(HideTipAction)).
			Post("/theme", new(ThemeAction)).

			EndAll()
	})
}
