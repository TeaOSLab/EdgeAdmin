package nodes

import (
	"errors"
	"fmt"
	"github.com/TeaOSLab/EdgeAdmin/internal/configloaders"
	"github.com/TeaOSLab/EdgeAdmin/internal/configs"
	teaconst "github.com/TeaOSLab/EdgeAdmin/internal/const"
	"github.com/TeaOSLab/EdgeAdmin/internal/events"
	"github.com/TeaOSLab/EdgeAdmin/internal/utils"
	"github.com/TeaOSLab/EdgeCommon/pkg/nodeconfigs"
	"github.com/iwind/TeaGo"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/rands"
	"github.com/iwind/TeaGo/types"
	"github.com/iwind/gosock/pkg/gosock"
	"gopkg.in/yaml.v3"
	"log"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

var SharedAdminNode *AdminNode = nil

type AdminNode struct {
	sock    *gosock.Sock
	subPIDs []int
}

func NewAdminNode() *AdminNode {
	return &AdminNode{}
}

func (this *AdminNode) Run() {
	SharedAdminNode = this

	// 启动管理界面
	var secret = this.genSecret()
	configs.Secret = secret

	// 本地Sock
	err := this.listenSock()
	if err != nil {
		logs.Println("[NODE]", err.Error())
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

	// 添加端口到防火墙
	this.addPortsToFirewall()

	// 监听信号
	var sigQueue = make(chan os.Signal, 8)
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

	// 设置DNS相关
	this.setupDNS()

	// 启动IP库
	this.startIPLibrary()

	// 启动Web服务
	sessionManager, err := NewSessionManager()
	if err != nil {
		log.Fatal("start session failed: " + err.Error())
		return
	}
	TeaGo.NewServer(false).
		AccessLog(false).
		EndAll().
		Session(sessionManager, teaconst.CookieSID).
		ReadHeaderTimeout(3 * time.Second).
		ReadTimeout(1200 * time.Second).
		Start()
}

// Daemon 实现守护进程
func (this *AdminNode) Daemon() {
	var sock = gosock.NewTmpSock(teaconst.ProcessName)
	var isDebug = lists.ContainsString(os.Args, "debug")
	for {
		conn, err := sock.Dial()
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

// InstallSystemService 安装系统服务
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

// AddSubPID 添加子PID
func (this *AdminNode) AddSubPID(pid int) {
	this.subPIDs = append(this.subPIDs, pid)
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
		data, err := os.ReadFile(templateFile)
		if err == nil {
			err = os.WriteFile(configFile, data, 0666)
			if err != nil {
				return fmt.Errorf("create config file failed: %w", err)
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
			err = os.WriteFile(configFile, []byte(templateYAML), 0666)
			if err != nil {
				return fmt.Errorf("create config file failed: %w", err)
			}
		}
	} else {
		return fmt.Errorf("can not read config from 'configs/server.yaml': %w", err)
	}

	return nil
}

// 添加端口到防火墙
func (this *AdminNode) addPortsToFirewall() {
	var configFile = Tea.ConfigFile("server.yaml")
	data, err := os.ReadFile(configFile)
	if err != nil {
		return
	}

	var config = &TeaGo.ServerConfig{}
	err = yaml.Unmarshal(data, config)
	if err != nil {
		return
	}

	var ports = []int{}
	if config.Http.On {
		for _, listen := range config.Http.Listen {
			_, portString, _ := net.SplitHostPort(listen)
			var port = types.Int(portString)
			if port > 0 && !lists.ContainsInt(ports, port) {
				ports = append(ports, port)
			}
		}
	}

	if config.Https.On {
		for _, listen := range config.Https.Listen {
			_, portString, _ := net.SplitHostPort(listen)
			var port = types.Int(portString)
			if port > 0 && !lists.ContainsInt(ports, port) {
				ports = append(ports, port)
			}
		}
	}

	utils.AddPortsToFirewall(ports)
}

// 设置DNS相关
func (this *AdminNode) setupDNS() {
	config, loadErr := configloaders.LoadAdminUIConfig()
	if loadErr != nil {
		// 默认使用go原生
		err := os.Setenv("GODEBUG", "netdns=go")
		if err != nil {
			logs.Println("[DNS_RESOLVER]set env failed: " + err.Error())
		}
		return
	}

	var err error
	switch config.DNSResolver.Type {
	case nodeconfigs.DNSResolverTypeGoNative:
		err = os.Setenv("GODEBUG", "netdns=go")
	case nodeconfigs.DNSResolverTypeCGO:
		err = os.Setenv("GODEBUG", "netdns=cgo")
	default:
		// 默认使用go原生
		err = os.Setenv("GODEBUG", "netdns=go")
	}
	if err != nil {
		logs.Println("[DNS_RESOLVER]set env failed: " + err.Error())
	}
}

// 启动API节点
func (this *AdminNode) startAPINode() {
	var configPath = Tea.Root + "/edge-api/configs/api.yaml"
	_, err := os.Stat(configPath)
	canStart := false
	if err == nil {
		canStart = true
	} else if err != nil && os.IsNotExist(err) {
		// 尝试恢复api.yaml
		homeDir, _ := os.UserHomeDir()
		paths := []string{}
		if len(homeDir) > 0 {
			paths = append(paths, homeDir+"/.edge-api/api.yaml")
		}
		paths = append(paths, "/etc/edge-api/api.yaml")
		for _, path := range paths {
			_, err = os.Stat(path)
			if err == nil {
				data, err := os.ReadFile(path)
				if err == nil {
					err = os.WriteFile(configPath, data, 0666)
					if err == nil {
						logs.Println("[NODE]recover 'edge-api/configs/api.yaml' from '" + path + "'")
						canStart = true
						break
					}
				}
			}
		}
	}

	dbPath := Tea.Root + "/edge-api/configs/db.yaml"
	_, err = os.Stat(dbPath)
	if err != nil && os.IsNotExist(err) {
		// 尝试恢复db.yaml
		homeDir, _ := os.UserHomeDir()
		paths := []string{}
		if len(homeDir) > 0 {
			paths = append(paths, homeDir+"/.edge-api/db.yaml")
		}
		paths = append(paths, "/etc/edge-api/db.yaml")
		for _, path := range paths {
			_, err = os.Stat(path)
			if err == nil {
				data, err := os.ReadFile(path)
				if err == nil {
					err = os.WriteFile(dbPath, data, 0666)
					if err == nil {
						logs.Println("[NODE]recover 'edge-api/configs/db.yaml' from '" + path + "'")
						break
					}
				}
			}
		}
	}

	if canStart {
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
	data, err := os.ReadFile(tmpFile)
	if err == nil && len(data) == 32 {
		return string(data)
	}
	secret := rands.String(32)
	_ = os.WriteFile(tmpFile, []byte(secret), 0666)
	return secret
}

// 监听本地sock
func (this *AdminNode) listenSock() error {
	this.sock = gosock.NewTmpSock(teaconst.ProcessName)

	// 检查是否在运行
	if this.sock.IsListening() {
		reply, err := this.sock.Send(&gosock.Command{Code: "pid"})
		if err == nil {
			return errors.New("error: the process is already running, pid: " + maps.NewMap(reply.Params).GetString("pid"))
		} else {
			return errors.New("error: the process is already running")
		}
	}

	// 启动监听
	go func() {
		this.sock.OnCommand(func(cmd *gosock.Command) {
			switch cmd.Code {
			case "pid":
				_ = cmd.Reply(&gosock.Command{
					Code: "pid",
					Params: map[string]interface{}{
						"pid": os.Getpid(),
					},
				})
			case "stop":
				_ = cmd.ReplyOk()

				// 关闭子进程
				for _, pid := range this.subPIDs {
					p, err := os.FindProcess(pid)
					if err == nil && p != nil {
						_ = p.Kill()
					}
				}

				// 停止当前目录下的API节点
				var apiSock = gosock.NewTmpSock("edge-api")
				apiReply, err := apiSock.Send(&gosock.Command{Code: "info"})
				if err == nil {
					adminExe, _ := os.Executable()
					if len(adminExe) > 0 && apiReply != nil && strings.HasPrefix(maps.NewMap(apiReply.Params).GetString("path"), filepath.Dir(filepath.Dir(adminExe))) {
						_, _ = apiSock.Send(&gosock.Command{Code: "stop"})
					}
				}

				// 退出主进程
				events.Notify(events.EventQuit)
				os.Exit(0)
			case "recover":
				teaconst.IsRecoverMode = true
				_ = cmd.ReplyOk()
			case "demo":
				teaconst.IsDemoMode = !teaconst.IsDemoMode
				_ = cmd.Reply(&gosock.Command{
					Params: map[string]interface{}{"isDemo": teaconst.IsDemoMode},
				})
			case "info":
				exePath, _ := os.Executable()
				_ = cmd.Reply(&gosock.Command{
					Code: "info",
					Params: map[string]interface{}{
						"pid":     os.Getpid(),
						"version": teaconst.Version,
						"path":    exePath,
					},
				})
			case "dev": // 切换到dev
				Tea.Env = Tea.EnvDev
				_ = cmd.ReplyOk()
			case "prod": // 切换到prod
				Tea.Env = Tea.EnvProd
				_ = cmd.ReplyOk()
			case "security.reset":
				var newConfig = configloaders.NewSecurityConfig()
				_ = configloaders.UpdateSecurityConfig(newConfig)
				_ = cmd.ReplyOk()
			}
		})

		err := this.sock.Listen()
		if err != nil {
			logs.Println("NODE", err.Error())
		}
	}()

	events.On(events.EventQuit, func() {
		logs.Println("NODE", "quit unix sock")
		_ = this.sock.Close()
	})

	return nil
}
