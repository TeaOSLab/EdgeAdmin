package main

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/apps"
	teaconst "github.com/TeaOSLab/EdgeAdmin/internal/const"
	_ "github.com/TeaOSLab/EdgeAdmin/internal/web"
	"github.com/iwind/TeaGo"
	_ "github.com/iwind/TeaGo/bootstrap"
	"github.com/iwind/TeaGo/rands"
	"github.com/iwind/TeaGo/sessions"
)

func main() {
	app := apps.NewAppCmd().
		Version(teaconst.Version).
		Product(teaconst.ProductName).
		Usage(teaconst.ProcessName + " [-v|start|stop|restart]")

	app.Run(func() {
		// 启动管理界面
		server := TeaGo.NewServer(false).
			AccessLog(false).
			EndAll().

			Session(sessions.NewFileSessionManager(
				86400,
				rands.String(32),
			))
		server.Start()
	})
}
