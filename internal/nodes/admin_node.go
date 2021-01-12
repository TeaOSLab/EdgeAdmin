package nodes

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/configs"
	teaconst "github.com/TeaOSLab/EdgeAdmin/internal/const"
	"github.com/TeaOSLab/EdgeAdmin/internal/errors"
	"github.com/TeaOSLab/EdgeAdmin/internal/events"
	"github.com/TeaOSLab/EdgeAdmin/internal/utils"
	"github.com/iwind/TeaGo"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/rands"
	"github.com/iwind/TeaGo/sessions"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"
)

type AdminNode struct {
	subPIDs []int
}

func NewAdminNode() *AdminNode {
	return &AdminNode{}
}

func (this *AdminNode) Run() {
	// 启动管理界面
	secret := this.genSecret()
	configs.Secret = secret

	// 本地Sock
	err := this.listenSock()
	if err != nil {
		logs.Println("NODE" + err.Error())
		return
	}

	// 检查server配置
	err = this.checkServer()
	if err != nil {
		if err != nil {
			logs.Println("[NODE]" + err.Error())
			return
		}
		return
	}

	// 监听信号
	sigQueue := make(chan os.Signal)
	signal.Notify(sigQueue, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL, syscall.SIGQUIT)
	go func() {
		for range sigQueue {
			for _, pid := range this.subPIDs {
				p, err := os.FindProcess(pid)
				if err == nil && p != nil {
					_ = p.Kill()
				}
			}
			os.Exit(0)
		}
	}()

	// 触发事件
	events.Notify(events.EventStart)

	// 启动API节点
	this.startAPINode()

	// 启动Web服务
	TeaGo.NewServer(false).
		AccessLog(false).
		EndAll().
		Session(sessions.NewFileSessionManager(86400, secret), teaconst.CookieSID).
		ReadHeaderTimeout(3 * time.Second).
		ReadTimeout(600 * time.Second).
		Start()
}

// 实现守护进程
func (this *AdminNode) Daemon() {
	path := os.TempDir() + "/edge-admin.sock"
	isDebug := lists.ContainsString(os.Args, "debug")
	isDebug = true
	for {
		conn, err := net.DialTimeout("unix", path, 1*time.Second)
		if err != nil {
			if isDebug {
				log.Println("[DAEMON]starting ...")
			}

			// 尝试启动
			err = func() error {
				exe, err := os.Executable()
				if err != nil {
					return err
				}
				cmd := exec.Command(exe)
				err = cmd.Start()
				if err != nil {
					return err
				}
				err = cmd.Wait()
				if err != nil {
					return err
				}
				return nil
			}()

			if err != nil {
				if isDebug {
					log.Println("[DAEMON]", err)
				}
				time.Sleep(1 * time.Second)
			} else {
				time.Sleep(5 * time.Second)
			}
		} else {
			_ = conn.Close()
			time.Sleep(5 * time.Second)
		}
	}
}

// 安装系统服务
func (this *AdminNode) InstallSystemService() error {
	shortName := teaconst.SystemdServiceName

	exe, err := os.Executable()
	if err != nil {
		return err
	}

	manager := utils.NewServiceManager(shortName, teaconst.ProductName)
	err = manager.Install(exe, []string{})
	if err != nil {
		return err
	}
	return nil
}

// 检查Server配置
func (this *AdminNode) checkServer() error {
	configFile := Tea.ConfigFile("server.yaml")
	_, err := os.Stat(configFile)
	if err == nil {
		return nil
	}

	if os.IsNotExist(err) {
		// 创建文件
		templateFile := Tea.ConfigFile("server.template.yaml")
		data, err := ioutil.ReadFile(templateFile)
		if err == nil {
			err = ioutil.WriteFile(configFile, data, 0666)
			if err != nil {
				return errors.New("create config file failed: " + err.Error())
			}
		} else {
			templateYAML := `# environment code
env: prod

# http
http:
  "on": true
  listen: [ "0.0.0.0:7788" ]

# https
https:
  "on": false
  listen: [ "0.0.0.0:443"]
  cert: ""
  key: ""
`
			err = ioutil.WriteFile(configFile, []byte(templateYAML), 0666)
			if err != nil {
				return errors.New("create config file failed: " + err.Error())
			}
		}
	} else {
		return errors.New("can not read config from 'configs/server.yaml': " + err.Error())
	}

	return nil
}

// 启动API节点
func (this *AdminNode) startAPINode() {
	_, err := os.Stat(Tea.Root + "/edge-api/configs/api.yaml")
	if err == nil {
		logs.Println("start edge-api")
		cmd := exec.Command(Tea.Root + "/edge-api/bin/edge-api")
		err = cmd.Start()
		if err != nil {
			logs.Println("[ERROR]start edge-api failed: " + err.Error())
		} else {
			this.subPIDs = append(this.subPIDs, cmd.Process.Pid)
		}
	}
}

// 生成Secret
func (this *AdminNode) genSecret() string {
	tmpFile := os.TempDir() + "/edge-admin-secret.tmp"
	data, err := ioutil.ReadFile(tmpFile)
	if err == nil && len(data) == 32 {
		return string(data)
	}
	secret := rands.String(32)
	_ = ioutil.WriteFile(tmpFile, []byte(secret), 0666)
	return secret
}

// 监听本地sock
func (this *AdminNode) listenSock() error {
	path := os.TempDir() + "/edge-admin.sock"

	// 检查是否已经存在
	_, err := os.Stat(path)
	if err == nil {
		conn, err := net.Dial("unix", path)
		if err != nil {
			_ = os.Remove(path)
		} else {
			_ = conn.Close()
		}
	}

	// 新的监听任务
	listener, err := net.Listen("unix", path)
	if err != nil {
		return err
	}
	events.On(events.EventQuit, func() {
		logs.Println("NODE", "quit unix sock")
		_ = listener.Close()
	})

	go func() {
		for {
			_, err := listener.Accept()
			if err != nil {
				return
			}
		}
	}()

	return nil
}
