package main

import (
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
		Usage(teaconst.ProcessName + " [-v|start|stop|restart]")

	app.Run(func() {
		adminNode := nodes.NewAdminNode()
		adminNode.Run()
	})
}
