package main

import (
	"fmt"
	"github.com/TeaOSLab/EdgeAdmin/internal/apps"
	"github.com/TeaOSLab/EdgeAdmin/internal/configs"
	teaconst "github.com/TeaOSLab/EdgeAdmin/internal/const"
	"github.com/TeaOSLab/EdgeAdmin/internal/nodes"
	_ "github.com/TeaOSLab/EdgeAdmin/internal/web"
	_ "github.com/iwind/TeaGo/bootstrap"
)

func main() {
	app := apps.NewAppCmd().
		Version(teaconst.Version).
		Product(teaconst.ProductName).
		Usage(teaconst.ProcessName+" [-v|start|stop|restart|service|daemon|reset]").
		Option("-h", "show this help").
		Option("-v", "show version").
		Option("start", "start the service").
		Option("stop", "stop the service").
		Option("service", "register service into systemd").
		Option("daemon", "start the service with daemon").
		Option("reset", "reset configs")

	app.On("daemon", func() {
		nodes.NewAdminNode().Daemon()
	})
	app.On("service", func() {
		err := nodes.NewAdminNode().InstallSystemService()
		if err != nil {
			fmt.Println("[ERROR]install failed: " + err.Error())
			return
		}
		fmt.Println("done")
	})
	app.On("reset", func() {
		err := configs.ResetAPIConfig()
		if err != nil {
			fmt.Println("[ERROR]reset failed: " + err.Error())
			return
		}
		fmt.Println("done")
	})
	app.Run(func() {
		adminNode := nodes.NewAdminNode()
		adminNode.Run()
	})
}
