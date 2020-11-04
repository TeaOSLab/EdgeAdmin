package iplibrary

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/servers/components/componentutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/helpers"
	"github.com/iwind/TeaGo"
)

func init() {
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		server.
			Helper(helpers.NewUserMustAuth()).
			Helper(NewHelper()).
			Helper(componentutils.NewComponentHelper()).
			Prefix("/servers/components/ip-library").
			Get("", new(IndexAction)).
			GetPost("/uploadPopup", new(UploadPopupAction)).
			Post("/delete", new(DeleteAction)).
			Get("/download", new(DownloadAction)).
			EndAll()
	})
}
