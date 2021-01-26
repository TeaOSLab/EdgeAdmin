package main

import (
	"fmt"
	"github.com/TeaOSLab/EdgeAdmin/internal/apps"
	teaconst "github.com/TeaOSLab/EdgeAdmin/internal/const"
	"github.com/TeaOSLab/EdgeAdmin/internal/nodes"
	_ "github.com/TeaOSLab/EdgeAdmin/internal/web"
	_ "github.com/iwind/TeaGo/bootstrap"
)

func main() {
	app := apps.NewAppCmd().
		Version(teaconst.Version).
		Product(teaconst.ProductName).
		Usage(teaconst.ProcessName + " [-v|start|stop|restart|service|daemon]")

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
	app.Run(func() {
		adminNode := nodes.NewAdminNode()
		adminNode.Run()
	})
}
