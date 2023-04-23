package certs

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/configloaders"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/servers/certs/acme"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/servers/certs/acme/accounts"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/servers/certs/acme/users"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/servers/certs/ocsp"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/helpers"
	"github.com/iwind/TeaGo"
)

func init() {
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		server.
			Helper(helpers.NewUserMustAuth(configloaders.AdminModuleCodeServer)).
			Helper(NewHelper()).
			Data("teaMenu", "servers").
			Data("teaSubMenu", "cert").
			Prefix("/servers/certs").
			Data("leftMenuItem", "cert").
			Get("", new(IndexAction)).
			GetPost("/uploadPopup", new(UploadPopupAction)).
			GetPost("/uploadBatchPopup", new(UploadBatchPopupAction)).
			Post("/delete", new(DeleteAction)).
			GetPost("/updatePopup", new(UpdatePopupAction)).
			Get("/certPopup", new(CertPopupAction)).
			Get("/viewKey", new(ViewKeyAction)).
			Get("/viewCert", new(ViewCertAction)).
			Get("/downloadKey", new(DownloadKeyAction)).
			Get("/downloadCert", new(DownloadCertAction)).
			Get("/downloadZip", new(DownloadZipAction)).
			Get("/selectPopup", new(SelectPopupAction)).
			Get("/datajs", new(DatajsAction)).

			// ACME任务
			Prefix("/servers/certs/acme").
			Data("leftMenuItem", "acme").
			Get("", new(acme.IndexAction)).
			GetPost("/create", new(acme.CreateAction)).
			Post("/run", new(acme.RunAction)).
			GetPost("/updateTaskPopup", new(acme.UpdateTaskPopupAction)).
			Post("/deleteTask", new(acme.DeleteTaskAction)).
			Post("/userOptions", new(acme.UserOptionsAction)).

			// ACME用户
			Prefix("/servers/certs/acme/users").
			Get("", new(users.IndexAction)).
			GetPost("/createPopup", new(users.CreatePopupAction)).
			GetPost("/updatePopup", new(users.UpdatePopupAction)).
			Post("/delete", new(users.DeleteAction)).
			GetPost("/selectPopup", new(users.SelectPopupAction)).
			Post("/accountsWithCode", new(users.AccountsWithCodeAction)).

			// ACME账号
			Prefix("/servers/certs/acme/accounts").
			Get("", new(accounts.IndexAction)).
			GetPost("/createPopup", new(accounts.CreatePopupAction)).
			GetPost("/updatePopup", new(accounts.UpdatePopupAction)).
			Post("/delete", new(accounts.DeleteAction)).

			// OCSP
			Prefix("/servers/certs/ocsp").
			Data("leftMenuItem", "ocsp").
			Get("", new(ocsp.IndexAction)).
			Post("/reset", new(ocsp.ResetAction)).
			Post("/resetAll", new(ocsp.ResetAllAction)).
			Post("/ignore", new(ocsp.IgnoreAction)).

			//
			EndAll()
	})
}
