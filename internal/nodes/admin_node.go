package nodes

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/errors"
	"github.com/iwind/TeaGo"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/rands"
	"github.com/iwind/TeaGo/sessions"
	"io/ioutil"
	"os"
	"os/exec"
)

type AdminNode struct {
}

func NewAdminNode() *AdminNode {
	return &AdminNode{}
}

func (this *AdminNode) Run() {
	// 启动管理界面
	secret := rands.String(32)

	// 测试环境下设置一个固定的key，方便我们调试
	if Tea.IsTesting() {
		secret = "8f983f4d69b83aaa0d74b21a212f6967"
	}

	// 检查server配置
	err := this.checkServer()
	if err != nil {
		return
	}

	// 启动API节点
	this.startAPINode()

	server := TeaGo.NewServer(false).
		AccessLog(false).
		EndAll().

		Session(sessions.NewFileSessionManager(86400, secret))
	server.Start()
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
func (this AdminNode) startAPINode() {
	_, err := os.Stat(Tea.Root + "/edge-api/configs/api.yaml")
	if err == nil {
		logs.Println("start edge-api")
		err = exec.Command(Tea.Root + "/edge-api/bin/edge-api").Start()
		if err != nil {
			logs.Println("[ERROR]start edge-api failed: " + err.Error())
		}
	}
}
