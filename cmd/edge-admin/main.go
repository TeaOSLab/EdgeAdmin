package main

import (
	"fmt"
	"github.com/TeaOSLab/EdgeAdmin/internal/apps"
	"github.com/TeaOSLab/EdgeAdmin/internal/configs"
	teaconst "github.com/TeaOSLab/EdgeAdmin/internal/const"
	"github.com/TeaOSLab/EdgeAdmin/internal/nodes"
	_ "github.com/TeaOSLab/EdgeAdmin/internal/web"
	_ "github.com/iwind/TeaGo/bootstrap"
	"github.com/iwind/gosock/pkg/gosock"
)

func main() {
	app := apps.NewAppCmd().
		Version(teaconst.Version).
		Product(teaconst.ProductName).
		Usage(teaconst.ProcessName+" [-v|start|stop|restart|service|daemon|reset|recover|demo]").
		Option("-h", "show this help").
		Option("-v", "show version").
		Option("start", "start the service").
		Option("stop", "stop the service").
		Option("service", "register service into systemd").
		Option("daemon", "start the service with daemon").
		Option("reset", "reset configs").
		Option("recover", "enter recovery mode")

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
	app.On("recover", func() {
		sock := gosock.NewTmpSock(teaconst.ProcessName)
		if !sock.IsListening() {
			fmt.Println("[ERROR]the service not started yet, you should start the service first")
			return
		}
		_, err := sock.Send(&gosock.Command{Code: "recover"})
		if err != nil {
			fmt.Println("[ERROR]enter recovery mode failed: " + err.Error())
			return
		}
		fmt.Println("enter recovery mode successfully")
	})
	app.On("demo", func() {
		sock := gosock.NewTmpSock(teaconst.ProcessName)
		if !sock.IsListening() {
			fmt.Println("[ERROR]the service not started yet, you should start the service first")
			return
		}
		_, err := sock.Send(&gosock.Command{Code: "demo"})
		if err != nil {
			fmt.Println("[ERROR]change demo mode failed: " + err.Error())
			return
		}
		fmt.Println("change demo mode successfully")
	})
	app.Run(func() {
		adminNode := nodes.NewAdminNode()
		adminNode.Run()
	})
}
