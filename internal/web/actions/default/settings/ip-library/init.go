package iplibrary

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
			Helper(NewHelper()).
			Helper(settingutils.NewHelper("ipLibrary")).
			Prefix("/settings/ip-library").
			Get("", new(IndexAction)).
			GetPost("/uploadPopup", new(UploadPopupAction)).
			Post("/delete", new(DeleteAction)).
			Get("/download", new(DownloadAction)).
			EndAll()
	})
}
