package main

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/apps"
	teaconst "github.com/TeaOSLab/EdgeAdmin/internal/const"
	_ "github.com/TeaOSLab/EdgeAdmin/internal/web"
	"github.com/iwind/TeaGo"
	"github.com/iwind/TeaGo/Tea"
	_ "github.com/iwind/TeaGo/bootstrap"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/rands"
	"github.com/iwind/TeaGo/sessions"
	"os"
	"os/exec"
)

func main() {
	app := apps.NewAppCmd().
		Version(teaconst.Version).
		Product(teaconst.ProductName).
		Usage(teaconst.ProcessName + " [-v|start|stop|restart]")

	app.Run(func() {
		// 启动管理界面
		secret := rands.String(32)

		// 测试环境下设置一个固定的key，方便我们调试
		if Tea.IsTesting() {
			secret = "8f983f4d69b83aaa0d74b21a212f6967"
		}

		// 启动API节点
		_, err := os.Stat(Tea.Root + "/edge-api/configs/api.yaml")
		if err == nil {
			logs.Println("start edge-api")
			err = exec.Command(Tea.Root + "/edge-api/bin/edge-api").Start()
			if err != nil {
				logs.Println("[ERROR]start edge-api failed: " + err.Error())
			}
		}

		server := TeaGo.NewServer(false).
			AccessLog(false).
			EndAll().

			Session(sessions.NewFileSessionManager(86400, secret))
		server.Start()
	})
}
