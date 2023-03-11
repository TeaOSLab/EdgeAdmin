package setup

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/setup/mysql"
	"github.com/iwind/TeaGo"
)

func init() {
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		server.
			Helper(new(Helper)).
			Prefix("/setup").
			Get("", new(IndexAction)).
			Post("/validateApi", new(ValidateApiAction)).
			Post("/validateDb", new(ValidateDbAction)).
			Post("/validateAdmin", new(ValidateAdminAction)).
			Post("/install", new(InstallAction)).
			Post("/status", new(StatusAction)).
			Post("/detectDB", new(DetectDBAction)).
			Post("/checkLocalIP", new(CheckLocalIPAction)).
			GetPost("/mysql/installPopup", new(mysql.InstallPopupAction)).
			Post("/mysql/installLogs", new(mysql.InstallLogsAction)).
			EndAll()
	})
}
