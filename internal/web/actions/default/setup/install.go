package setup

import (
	"bytes"
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/configs"
	"github.com/TeaOSLab/EdgeAdmin/internal/nodes"
	"github.com/TeaOSLab/EdgeAdmin/internal/rpc"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/configutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/systemconfigs"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/gosock/pkg/gosock"
	"gopkg.in/yaml.v3"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"
)

type InstallAction struct {
	actionutils.ParentAction

	apiSetupFinished bool
}

func (this *InstallAction) RunPost(params struct {
	ApiNodeJSON []byte
	DbJSON      []byte
	AdminJSON   []byte

	Must *actions.Must
}) {
	currentStatusText = ""
	defer func() {
		currentStatusText = ""
	}()

	// API节点配置
	currentStatusText = "正在检查API节点配置"
	var apiNodeMap = maps.Map{}
	err := json.Unmarshal(params.ApiNodeJSON, &apiNodeMap)
	if err != nil {
		this.Fail("API节点配置数据解析错误，请刷新页面后重新尝试安装，错误信息：" + err.Error())
	}

	// 数据库
	currentStatusText = "正在检查数据库配置"
	var dbMap = maps.Map{}
	err = json.Unmarshal(params.DbJSON, &dbMap)
	if err != nil {
		this.Fail("数据库配置数据解析错误，请刷新页面后重新尝试安装，错误信息：" + err.Error())
	}

	// 管理员
	currentStatusText = "正在检查管理员配置"
	var adminMap = maps.Map{}
	err = json.Unmarshal(params.AdminJSON, &adminMap)
	if err != nil {
		this.Fail("管理员数据解析错误，请刷新页面后重新尝试安装，错误信息：" + err.Error())
	}

	// 安装API节点
	var mode = apiNodeMap.GetString("mode")
	if mode == "new" {
		currentStatusText = "准备启动新API节点"

		// 整个系统目录结构为：
		//  edge-admin/
		//    edge-api/
		//    bin/
		//    ...

		// 检查环境
		var apiNodeDir = Tea.Root + "/edge-api"
		for _, dir := range []string{"edge-api", "edge-api/configs", "edge-api/bin"} {
			var searchDir = Tea.Root + "/" + dir
			_, err = os.Stat(searchDir)
			if err != nil {
				if os.IsNotExist(err) {
					this.Fail("在当前目录（" + Tea.Root + "）下找不到" + dir + "目录，请将" + dir + "目录上传或者重新下载解压")
				}
				this.Fail("无法检查" + dir + "目录，发生错误：" + err.Error())
				return
			}
		}

		// 保存数据库配置
		var dbConfig = &configs.SimpleDBConfig{
			User:     dbMap.GetString("username"),
			Password: dbMap.GetString("password"),
			Database: dbMap.GetString("database"),
			Host:     configutils.QuoteIP(dbMap.GetString("host")) + ":" + dbMap.GetString("port"),
		}
		dbConfigData, err := yaml.Marshal(dbConfig)
		if err != nil {
			this.Fail("生成数据库配置失败：" + err.Error())
			return
		}
		err = os.WriteFile(apiNodeDir+"/configs/db.yaml", dbConfigData, 0666)
		if err != nil {
			this.Fail("保存数据库配置失败（db.yaml）：" + err.Error())
			return
		}

		err = dbConfig.GenerateOldConfig(apiNodeDir + "/configs/.db.yaml")
		if err != nil {
			this.Fail("保存数据库配置失败（.db.yaml）：" + err.Error())
			return
		}

		// 生成备份文件
		homeDir, _ := os.UserHomeDir()
		var backupDirs = []string{"/etc/edge-api"}
		if len(homeDir) > 0 {
			backupDirs = append(backupDirs, homeDir+"/.edge-api")
		}
		for _, backupDir := range backupDirs {
			stat, err := os.Stat(backupDir)
			if err == nil && stat.IsDir() {
				_ = os.WriteFile(backupDir+"/db.yaml", dbConfigData, 0666)
			} else if err != nil && os.IsNotExist(err) {
				err = os.Mkdir(backupDir, 0777)
				if err == nil {
					_ = os.WriteFile(backupDir+"/db.yaml", dbConfigData, 0666)
				}
			}
		}

		err = os.WriteFile(Tea.ConfigFile("/api_db.yaml"), dbConfigData, 0666)
		if err != nil {
			this.Fail("保存数据库配置失败（api_db.yaml）：" + err.Error())
			return
		}

		// 生成备份文件
		backupDirs = []string{"/etc/edge-admin"}
		if len(homeDir) > 0 {
			backupDirs = append(backupDirs, homeDir+"/.edge-admin")
		}
		for _, backupDir := range backupDirs {
			stat, err := os.Stat(backupDir)
			if err == nil && stat.IsDir() {
				_ = os.WriteFile(backupDir+"/api_db.yaml", dbConfigData, 0666)
			} else if err != nil && os.IsNotExist(err) {
				err = os.Mkdir(backupDir, 0777)
				if err == nil {
					_ = os.WriteFile(backupDir+"/api_db.yaml", dbConfigData, 0666)
				}
			}
		}

		// 开始安装
		currentStatusText = "正在安装数据库表结构并写入数据"
		var resultMap = maps.Map{}
		logs.Println("[INSTALL]setup edge-api")
		{
			this.apiSetupFinished = false
			var cmd = exec.Command(apiNodeDir+"/bin/edge-api", "setup", "-api-node-protocol=http", "-api-node-host=\""+apiNodeMap.GetString("newHost")+"\"", "-api-node-port=\""+apiNodeMap.GetString("newPort")+"\"")
			var output = bytes.NewBuffer(nil)
			cmd.Stdout = output

			var stderr = bytes.NewBuffer(nil)
			cmd.Stderr = stderr

			// 试图读取执行日志
			go this.startReadingAPIInstallLog()

			err = cmd.Run()
			this.apiSetupFinished = true
			if err != nil {
				this.Fail("安装失败：" + err.Error() + ": " + string(append(output.Bytes(), stderr.Bytes()...)))
			}

			var resultData = output.Bytes()
			err = json.Unmarshal(resultData, &resultMap)
			if err != nil {

				this.Fail("安装节点时返回数据错误：" + err.Error() + "(" + string(resultData) + ")")
			}
			if !resultMap.GetBool("isOk") {
				this.Fail("节点安装错误：" + resultMap.GetString("error"))
			}

			// 等数据完全写入
			time.Sleep(1 * time.Second)
		}

		// 关闭正在运行的API节点，防止冲突
		logs.Println("[INSTALL]stop edge-api")
		{
			var cmd = exec.Command(apiNodeDir+"/bin/edge-api", "stop")
			_ = cmd.Run()
		}

		// 启动API节点
		currentStatusText = "正在启动API节点"
		logs.Println("[INSTALL]start edge-api")
		{
			var cmd = exec.Command(apiNodeDir + "/bin/edge-api")
			err = cmd.Start()
			if err != nil {
				this.Fail("API节点启动失败：" + err.Error())
			}

			// 记录子PID方便退出的时候一起退出
			nodes.SharedAdminNode.AddSubPID(cmd.Process.Pid)

			// 等待API节点初始化完成
			currentStatusText = "正在等待API节点启动完毕"
			var apiNodeSock = gosock.NewTmpSock("edge-api")
			var maxRetries = 5
			for {
				reply, err := apiNodeSock.SendTimeout(&gosock.Command{
					Code: "starting",
				}, 3*time.Second)
				if err != nil {
					if maxRetries < 0 {
						this.Fail("API节点启动失败，请查看运行日志检查是否正常")
					} else {
						time.Sleep(3 * time.Second)
						maxRetries--
					}
				} else {
					if !maps.NewMap(reply.Params).GetBool("isStarting") {
						currentStatusText = "API节点启动完毕"
						break
					}

					// 继续等待完成
					time.Sleep(3 * time.Second)
				}
			}
		}

		// 写入API节点配置，完成安装
		var apiConfig = &configs.APIConfig{
			RPCEndpoints: []string{"http://" + configutils.QuoteIP(apiNodeMap.GetString("newHost")) + ":" + apiNodeMap.GetString("newPort")},
			NodeId:       resultMap.GetString("adminNodeId"),
			Secret:       resultMap.GetString("adminNodeSecret"),
		}

		// 设置管理员
		currentStatusText = "正在设置管理员"
		client, err := rpc.NewRPCClient(apiConfig, false)
		if err != nil {
			this.FailField("oldHost", "测试API节点时出错，请检查配置，错误信息："+err.Error())
		}
		ctx := client.Context(0)
		for i := 0; i < 3; i++ {
			_, err = client.AdminRPC().CreateOrUpdateAdmin(ctx, &pb.CreateOrUpdateAdminRequest{
				Username: adminMap.GetString("username"),
				Password: adminMap.GetString("password"),
			})
			// 这里我们尝试多次是为了等待API节点启动完毕
			if err != nil {
				time.Sleep(1 * time.Second)
			} else {
				break
			}
		}
		if err != nil {
			this.Fail("设置管理员账号出错：" + err.Error())
		}

		// 设置访问日志保留天数
		currentStatusText = "正在配置访问日志保留天数"
		var accessLogKeepDays = dbMap.GetInt("accessLogKeepDays")
		if accessLogKeepDays > 0 {
			var config = systemconfigs.NewDatabaseConfig()
			config.ServerAccessLog.Clean.Days = accessLogKeepDays
			configJSON, err := json.Marshal(config)
			if err != nil {
				this.Fail("配置设置访问日志保留天数出错：" + err.Error())
				return
			}
			_, err = client.SysSettingRPC().UpdateSysSetting(ctx, &pb.UpdateSysSettingRequest{
				Code:      systemconfigs.SettingCodeDatabaseConfigSetting,
				ValueJSON: configJSON,
			})
			if err != nil {
				this.Fail("配置设置访问日志保留天数出错：" + err.Error())
				return
			}
		}

		err = apiConfig.WriteFile(Tea.ConfigFile(configs.ConfigFileName))
		if err != nil {
			this.Fail("保存配置失败，原因：" + err.Error())
		}

		this.Success()
	} else if mode == "old" {
		// 构造RPC
		var apiConfig = &configs.APIConfig{
			RPCEndpoints: []string{apiNodeMap.GetString("oldProtocol") + "://" + configutils.QuoteIP(apiNodeMap.GetString("oldHost")) + ":" + apiNodeMap.GetString("oldPort")},
			NodeId:       apiNodeMap.GetString("oldNodeId"),
			Secret:       apiNodeMap.GetString("oldNodeSecret"),
		}
		client, err := rpc.NewRPCClient(apiConfig, false)
		if err != nil {
			this.FailField("oldHost", "测试API节点时出错，请检查配置，错误信息："+err.Error())
		}

		defer func() {
			_ = client.Close()
		}()

		// 设置管理员
		var ctx = client.APIContext(0)
		_, err = client.AdminRPC().CreateOrUpdateAdmin(ctx, &pb.CreateOrUpdateAdminRequest{
			Username: adminMap.GetString("username"),
			Password: adminMap.GetString("password"),
		})
		if err != nil {
			this.Fail("设置管理员账号出错：" + err.Error())
		}

		// 设置访问日志保留天数
		var accessLogKeepDays = dbMap.GetInt("accessLogKeepDays")
		if accessLogKeepDays > 0 {
			var config = systemconfigs.NewDatabaseConfig()
			config.ServerAccessLog.Clean.Days = accessLogKeepDays
			configJSON, err := json.Marshal(config)
			if err != nil {
				this.Fail("配置设置访问日志保留天数出错：" + err.Error())
				return
			}
			_, err = client.SysSettingRPC().UpdateSysSetting(ctx, &pb.UpdateSysSettingRequest{
				Code:      systemconfigs.SettingCodeDatabaseConfigSetting,
				ValueJSON: configJSON,
			})
			if err != nil {
				this.Fail("配置设置访问日志保留天数出错：" + err.Error())
				return
			}
		}

		// 写入API节点配置，完成安装
		err = apiConfig.WriteFile(Tea.ConfigFile(configs.ConfigFileName))
		if err != nil {
			this.Fail("保存配置失败，原因：" + err.Error())
		}

		// 成功
		this.Success()
	} else {
		this.Fail("错误的API节点模式：'" + mode + "'")
	}
}

// 读取API安装时的日志，以便于显示当前正在执行的任务
func (this *InstallAction) startReadingAPIInstallLog() {
	var tmpDir = os.TempDir()
	if len(tmpDir) == 0 {
		return
	}
	var logFile = tmpDir + "/edge-install.log"

	var logFp *os.File
	var err error

	// 尝试5秒钟
	for i := 0; i < 10; i++ {
		logFp, err = os.Open(logFile)
		if err != nil {
			time.Sleep(1 * time.Second)
			continue
		} else {
			break
		}
	}
	if err != nil {
		return
	}

	if this.apiSetupFinished {
		_ = logFp.Close()
		return
	}

	go func() {
		defer func() {
			_ = logFp.Close()
		}()

		var ticker = time.NewTicker(1 * time.Second)
		var logBuf = make([]byte, 256)
		for range ticker.C {
			if this.apiSetupFinished {
				return
			}

			_, err = logFp.Seek(-256, io.SeekEnd)
			if err != nil {
				currentStatusText = ""
				return
			}

			n, err := logFp.Read(logBuf)
			if err != nil {
				currentStatusText = ""
				return
			}

			if n > 0 {
				var logData = string(logBuf[:n])
				var lines = strings.Split(logData, "\n")
				if len(lines) >= 3 {
					var line = strings.TrimSpace(lines[len(lines)-2])
					if len(line) > 0 {
						if !this.apiSetupFinished {
							currentStatusText = "正在执行 " + line + " ..."
						}
					}
				}
			}
		}
	}()
}
