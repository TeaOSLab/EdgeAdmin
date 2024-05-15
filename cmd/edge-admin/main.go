package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/TeaOSLab/EdgeAdmin/internal/apps"
	"github.com/TeaOSLab/EdgeAdmin/internal/configs"
	teaconst "github.com/TeaOSLab/EdgeAdmin/internal/const"
	"github.com/TeaOSLab/EdgeAdmin/internal/gen"
	"github.com/TeaOSLab/EdgeAdmin/internal/nodes"
	"github.com/TeaOSLab/EdgeAdmin/internal/rpc"
	"github.com/TeaOSLab/EdgeAdmin/internal/utils"
	executils "github.com/TeaOSLab/EdgeAdmin/internal/utils/exec"
	_ "github.com/TeaOSLab/EdgeAdmin/internal/web"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/settings/updates/updateutils"
	_ "github.com/TeaOSLab/EdgeCommon/pkg/langs/messages"
	"github.com/iwind/TeaGo/Tea"
	_ "github.com/iwind/TeaGo/bootstrap"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/gosock/pkg/gosock"
	"log"
	"os"
	"os/exec"
	"time"
)

func main() {
	var app = apps.NewAppCmd().
		Version(teaconst.Version).
		Product(teaconst.ProductName).
		Usage(teaconst.ProcessName+" [-h|-v|start|stop|restart|service|daemon|reset|recover|demo|upgrade]").
		Usage(teaconst.ProcessName+" [dev|prod]").
		Option("-h", "show this help").
		Option("-v", "show version").
		Option("start", "start the service").
		Option("stop", "stop the service").
		Option("restart", "restart the service").
		Option("service", "register service into systemd").
		Option("daemon", "start the service with daemon").
		Option("reset", "reset configs").
		Option("recover", "enter recovery mode").
		Option("demo", "switch to demo mode").
		Option("dev", "switch to 'dev' mode").
		Option("prod", "switch to 'prod' mode").
		Option("upgrade [--url=URL]", "upgrade from official site or an url").
		Option("install-local-node", "install a local node").
		Option("security.reset", "reset security config")

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

		// reset local api
		var apiNodeExe = Tea.Root + "/edge-api/bin/edge-api"
		_, err = os.Stat(apiNodeExe)
		if err == nil {
			var cmd = exec.Command(apiNodeExe, "reset")
			var stderr = &bytes.Buffer{}
			cmd.Stderr = stderr
			err = cmd.Run()
			if err != nil {
				fmt.Println("reset api node failed: " + stderr.String())
			}
		}

		fmt.Println("done")
	})
	app.On("recover", func() {
		var sock = gosock.NewTmpSock(teaconst.ProcessName)
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
		var sock = gosock.NewTmpSock(teaconst.ProcessName)
		if !sock.IsListening() {
			fmt.Println("[ERROR]the service not started yet, you should start the service first")
			return
		}
		reply, err := sock.Send(&gosock.Command{Code: "demo"})
		if err != nil {
			fmt.Println("[ERROR]change demo mode failed: " + err.Error())
			return
		}
		var isDemo = maps.NewMap(reply.Params).GetBool("isDemo")
		if isDemo {
			fmt.Println("change demo mode to: on")
		} else {
			fmt.Println("change demo mode to: off")
		}
	})
	app.On("generate", func() {
		err := gen.Generate()
		if err != nil {
			fmt.Println("generate failed: " + err.Error())
			return
		}
	})
	app.On("dev", func() {
		var env = "dev"
		var sock = gosock.NewTmpSock(teaconst.ProcessName)
		_, err := sock.Send(&gosock.Command{
			Code:   env,
			Params: nil,
		})
		if err != nil {
			fmt.Println("failed to switch to '" + env + "': " + err.Error())
		} else {
			fmt.Println("switch to '" + env + "' ok")
		}
	})
	app.On("prod", func() {
		var env = "prod"
		var sock = gosock.NewTmpSock(teaconst.ProcessName)
		_, err := sock.Send(&gosock.Command{
			Code:   env,
			Params: nil,
		})
		if err != nil {
			fmt.Println("failed to switch to '" + env + "': " + err.Error())
		} else {
			fmt.Println("switch to '" + env + "' ok")
		}
	})
	app.On("upgrade", func() {
		var downloadURL = ""
		var flagSet = flag.NewFlagSet("", flag.ContinueOnError)
		flagSet.StringVar(&downloadURL, "url", "", "new version download url")
		_ = flagSet.Parse(os.Args[2:])

		var manager = utils.NewUpgradeManager("admin", downloadURL)
		log.Println("checking latest version ...")
		var ticker = time.NewTicker(1 * time.Second)
		go func() {
			var lastProgress float32 = 0
			var isStarted = false
			for range ticker.C {
				if manager.IsDownloading() {
					if !isStarted {
						if len(manager.NewVersion()) == 0 {
							continue
						}
						log.Println("start downloading v" + manager.NewVersion() + " ...")
						isStarted = true
					}
					var progress = manager.Progress()
					if progress >= 0 {
						if progress == 0 || progress == 1 || progress-lastProgress >= 0.1 {
							lastProgress = progress
							log.Printf("%.2f%%", manager.Progress()*100)
						}
					}
				} else {
					break
				}
			}
		}()
		err := manager.Start()
		if err != nil {
			log.Println("upgrade failed: " + err.Error())
			return
		}

		// try to exec local 'edge-api upgrade'
		rpcClient, err := rpc.SharedRPC()
		if err == nil {
			exePath, ok := updateutils.CheckLocalAPINode(rpcClient, rpcClient.Context(0))
			if ok && len(exePath) > 0 {
				log.Println("upgrading database ...")
				var cmd = executils.NewCmd(exePath, "upgrade")
				_ = cmd.Run()
			}
		}

		log.Println("finished!")
		log.Println("restarting ...")
		app.RunRestart()
	})
	app.On("security.reset", func() {
		var sock = gosock.NewTmpSock(teaconst.ProcessName)
		if !sock.IsListening() {
			fmt.Println("[ERROR]the service not started yet, you should start the service first")
			return
		}
		_, _ = sock.Send(&gosock.Command{
			Code: "security.reset",
		})
		fmt.Println("ok")
	})
	app.Run(func() {
		var adminNode = nodes.NewAdminNode()
		adminNode.Run()
	})
}
