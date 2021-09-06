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
	"github.com/go-yaml/yaml"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/dbs"
	"github.com/iwind/TeaGo/maps"
	"io/ioutil"
	"os"
	"os/exec"
	"time"
)

type InstallAction struct {
	actionutils.ParentAction
}

func (this *InstallAction) RunPost(params struct {
	ApiNodeJSON []byte
	DbJSON      []byte
	AdminJSON   []byte

	Must *actions.Must
}) {
	// API节点配置
	apiNodeMap := maps.Map{}
	err := json.Unmarshal(params.ApiNodeJSON, &apiNodeMap)
	if err != nil {
		this.Fail("API节点配置数据解析错误，请刷新页面后重新尝试安装，错误信息：" + err.Error())
	}

	// 数据库
	dbMap := maps.Map{}
	err = json.Unmarshal(params.DbJSON, &dbMap)
	if err != nil {
		this.Fail("数据库配置数据解析错误，请刷新页面后重新尝试安装，错误信息：" + err.Error())
	}

	// 管理员
	adminMap := maps.Map{}
	err = json.Unmarshal(params.AdminJSON, &adminMap)
	if err != nil {
		this.Fail("管理员数据解析错误，请刷新页面后重新尝试安装，错误信息：" + err.Error())
	}

	// 安装API节点
	mode := apiNodeMap.GetString("mode")
	if mode == "new" {
		// 整个系统目录结构为：
		//  edge-admin/
		//    edge-api/
		//    bin/
		//    ...

		// 检查环境
		apiNodeDir := Tea.Root + "/edge-api"
		for _, dir := range []string{"edge-api", "edge-api/configs", "edge-api/bin"} {
			apiNodeDir := Tea.Root + "/" + dir
			_, err = os.Stat(apiNodeDir)
			if err != nil {
				if os.IsNotExist(err) {
					this.Fail("在当前目录（" + Tea.Root + "）下找不到" + dir + "目录，请将" + dir + "目录上传或者重新下载解压")
				}
				this.Fail("无法检查" + dir + "目录，发生错误：" + err.Error())
			}
		}

		// 保存数据库配置
		dsn := dbMap.GetString("username") + ":" + dbMap.GetString("password") + "@tcp(" + configutils.QuoteIP(dbMap.GetString("host")) + ":" + dbMap.GetString("port") + ")/" + dbMap.GetString("database") + "?charset=utf8mb4&timeout=30s"
		dbConfig := &dbs.Config{
			DBs: map[string]*dbs.DBConfig{
				"prod": {
					Driver: "mysql",
					Dsn:    dsn,
					Prefix: "edge",
				}},
		}
		dbConfig.Default.DB = "prod"
		dbConfigData, err := yaml.Marshal(dbConfig)
		if err != nil {
			this.Fail("生成数据库配置失败：" + err.Error())
		}
		err = ioutil.WriteFile(apiNodeDir+"/configs/db.yaml", dbConfigData, 0666)
		if err != nil {
			this.Fail("保存数据库配置失败：" + err.Error())
		}

		// 生成备份文件
		homeDir, _ := os.UserHomeDir()
		backupDirs := []string{"/etc/edge-api"}
		if len(homeDir) > 0 {
			backupDirs = append(backupDirs, homeDir+"/.edge-api")
		}
		for _, backupDir := range backupDirs {
			stat, err := os.Stat(backupDir)
			if err == nil && stat.IsDir() {
				_ = ioutil.WriteFile(backupDir+"/db.yaml", dbConfigData, 0666)
			} else if err != nil && os.IsNotExist(err) {
				err = os.Mkdir(backupDir, 0777)
				if err == nil {
					_ = ioutil.WriteFile(backupDir+"/db.yaml", dbConfigData, 0666)
				}
			}
		}

		err = ioutil.WriteFile(Tea.ConfigFile("/api_db.yaml"), dbConfigData, 0666)
		if err != nil {
			this.Fail("保存数据库配置失败：" + err.Error())
		}

		// 生成备份文件
		backupDirs = []string{"/etc/edge-admin"}
		if len(homeDir) > 0 {
			backupDirs = append(backupDirs, homeDir+"/.edge-admin")
		}
		for _, backupDir := range backupDirs {
			stat, err := os.Stat(backupDir)
			if err == nil && stat.IsDir() {
				_ = ioutil.WriteFile(backupDir+"/api_db.yaml", dbConfigData, 0666)
			} else if err != nil && os.IsNotExist(err) {
				err = os.Mkdir(backupDir, 0777)
				if err == nil {
					_ = ioutil.WriteFile(backupDir+"/api_db.yaml", dbConfigData, 0666)
				}
			}
		}

		// 开始安装
		var resultMap = maps.Map{}
		{
			cmd := exec.Command(apiNodeDir+"/bin/edge-api", "setup", "-api-node-protocol=http", "-api-node-host=\""+apiNodeMap.GetString("newHost")+"\"", "-api-node-port=\""+apiNodeMap.GetString("newPort")+"\"")
			output := bytes.NewBuffer([]byte{})
			cmd.Stdout = output
			err = cmd.Run()
			if err != nil {
				this.Fail("安装失败：" + err.Error())
			}

			resultData := output.Bytes()
			err = json.Unmarshal(resultData, &resultMap)
			if err != nil {
				this.Fail("安装节点时返回数据错误：" + err.Error() + "(" + string(resultData) + ")")
			}
			if !resultMap.GetBool("isOk") {
				this.Fail("节点安装错误：" + resultMap.GetString("error"))
			}
		}

		// 启动API节点
		{
			cmd := exec.Command(apiNodeDir + "/bin/edge-api")
			err = cmd.Start()
			if err != nil {
				this.Fail("API节点启动失败：" + err.Error())
			}

			// 记录子PID方便退出的时候一起退出
			nodes.SharedAdminNode.AddSubPID(cmd.Process.Pid)

			// 等待API节点初始化完成
			time.Sleep(5 * time.Second)
		}

		// 写入API节点配置，完成安装
		apiConfig := &configs.APIConfig{
			RPC: struct {
				Endpoints []string `yaml:"endpoints"`
			}{
				Endpoints: []string{"http://" + configutils.QuoteIP(apiNodeMap.GetString("newHost")) + ":" + apiNodeMap.GetString("newPort")},
			},
			NodeId: resultMap.GetString("adminNodeId"),
			Secret: resultMap.GetString("adminNodeSecret"),
		}

		// 设置管理员
		client, err := rpc.NewRPCClient(apiConfig)
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
			}
		}
		if err != nil {
			this.Fail("设置管理员账号出错：" + err.Error())
		}

		// 设置访问日志保留天数
		var accessLogKeepDays = dbMap.GetInt("accessLogKeepDays")
		if accessLogKeepDays > 0 {
			var config = &systemconfigs.DatabaseConfig{}
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

		err = apiConfig.WriteFile(Tea.ConfigFile("api.yaml"))
		if err != nil {
			this.Fail("保存配置失败，原因：" + err.Error())
		}

		this.Success()
	} else if mode == "old" {
		// 构造RPC
		apiConfig := &configs.APIConfig{
			RPC: struct {
				Endpoints []string `yaml:"endpoints"`
			}{
				Endpoints: []string{apiNodeMap.GetString("oldProtocol") + "://" + configutils.QuoteIP(apiNodeMap.GetString("oldHost")) + ":" + apiNodeMap.GetString("oldPort")},
			},
			NodeId: apiNodeMap.GetString("oldNodeId"),
			Secret: apiNodeMap.GetString("oldNodeSecret"),
		}
		client, err := rpc.NewRPCClient(apiConfig)
		if err != nil {
			this.FailField("oldHost", "测试API节点时出错，请检查配置，错误信息："+err.Error())
		}

		// 设置管理员
		ctx := client.APIContext(0)
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
			var config = &systemconfigs.DatabaseConfig{}
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
		err = apiConfig.WriteFile(Tea.ConfigFile("api.yaml"))
		if err != nil {
			this.Fail("保存配置失败，原因：" + err.Error())
		}

		// 成功
		this.Success()
	} else {
		this.Fail("错误的API节点模式：'" + mode + "'")
	}
}
