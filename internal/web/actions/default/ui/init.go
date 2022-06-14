package ui

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/configloaders"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/helpers"
	"github.com/iwind/TeaGo"
	"github.com/iwind/TeaGo/Tea"
)

func init() {
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		server.
			Prefix("/ui").

			// 公共可以访问的链接
			Get("/image/:fileId", new(ImageAction)).

			// 以下需要登录
			Helper(helpers.NewUserMustAuth(configloaders.AdminModuleCodeCommon)).
			Get("/download", new(DownloadAction)).
			GetPost("/selectProvincesPopup", new(SelectProvincesPopupAction)).
			GetPost("/selectCountriesPopup", new(SelectCountriesPopupAction)).
			Post("/eventLevelOptions", new(EventLevelOptionsAction)).
			Post("/showTip", new(ShowTipAction)).
			Post("/hideTip", new(HideTipAction)).
			Post("/theme", new(ThemeAction)).
			Post("/validateIPs", new(ValidateIPsAction)).
			Post("/providerOptions", new(ProviderOptionsAction)).
			Post("/countryOptions", new(CountryOptionsAction)).
			Post("/provinceOptions", new(ProvinceOptionsAction)).
			Post("/cityOptions", new(CityOptionsAction)).
			EndAll()

		// 开发环境下总是动态加载，以便于调试
		if Tea.IsTesting() {
			server.
				Get("/js/components.js", new(ComponentsAction)).
				EndAll()
		}
	})
}
