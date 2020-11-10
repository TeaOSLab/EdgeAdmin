package nodes

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/configs"
	"github.com/TeaOSLab/EdgeAdmin/internal/errors"
	"github.com/iwind/TeaGo"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/rands"
	"github.com/iwind/TeaGo/sessions"
	"io/ioutil"
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

	// 检查server配置
	err := this.checkServer()
	if err != nil {
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

	// 启动API节点
	this.startAPINode()

	TeaGo.NewServer(false).
		AccessLog(false).
		EndAll().
		Session(sessions.NewFileSessionManager(86400, secret)).
		ReadHeaderTimeout(3 * time.Second).
		ReadTimeout(600 * time.Second).
		Start()
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
	tmpFile := os.TempDir() + "/edge-secret.tmp"
	data, err := ioutil.ReadFile(tmpFile)
	if err == nil && len(data) == 32 {
		return string(data)
	}
	secret := rands.String(32)
	_ = ioutil.WriteFile(tmpFile, []byte(secret), 0666)
	return secret
}
