package database

import (
	"fmt"
	"github.com/TeaOSLab/EdgeAdmin/internal/configs"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/configutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
	"github.com/go-sql-driver/mysql"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/dbs"
	"github.com/iwind/TeaGo/maps"
	"gopkg.in/yaml.v3"
	"net"
	"os"
	"regexp"
	"strings"
)

type UpdateAction struct {
	actionutils.ParentAction
}

func (this *UpdateAction) Init() {
	this.Nav("", "", "update")
}

func (this *UpdateAction) RunGet(params struct{}) {
	this.Data["dbConfig"] = maps.Map{
		"host":     "",
		"port":     "",
		"username": "",
		"password": "",
		"database": "",
	}

	configFile := Tea.ConfigFile("api_db.yaml")
	data, err := os.ReadFile(configFile)
	if err != nil {
		return
	}

	// new config
	var config = &configs.SimpleDBConfig{}
	err = yaml.Unmarshal(data, config)
	if err == nil && len(config.Host) > 0 {
		host, port, splitErr := net.SplitHostPort(config.Host)
		if splitErr != nil {
			port = "3306"
		}

		this.Data["dbConfig"] = maps.Map{
			"host":     host,
			"port":     port,
			"username": config.User,
			"password": config.Password,
			"database": config.Database,
		}

		this.Show()
		return
	}

	this.parseOldConfig(data)
}

func (this *UpdateAction) RunPost(params struct {
	Host     string
	Port     int32
	Database string
	Username string
	Password string

	Must *actions.Must
}) {
	defer this.CreateLogInfo(codes.Database_LogUpdateAPINodeDatabaseConfig)

	params.Must.
		Field("host", params.Host).
		Require("请输入主机地址").
		Expect(func() (message string, success bool) {
			// 是否为IP
			if net.ParseIP(params.Host) != nil {
				success = true
				return
			}
			if !regexp.MustCompile(`^[\w.-]+$`).MatchString(params.Host) {
				message = "主机地址中不能包含特殊字符"
				success = false
				return
			}
			success = true
			return
		}).
		Field("port", params.Port).
		Gt(0, "端口需要大于0").
		Lt(65535, "端口需要小于65535").
		Field("database", params.Database).
		Require("请输入数据库名称").
		Match(`^[\w\.-]+$`, "数据库名称中不能包含特殊字符").
		Field("username", params.Username).
		Require("请输入连接数据库的用户名").
		Match(`^[\w\.-]+$`, "用户名中不能包含特殊字符")

	var config = &configs.SimpleDBConfig{
		User:     params.Username,
		Password: params.Password,
		Database: params.Database,
		Host:     configutils.QuoteIP(params.Host) + ":" + fmt.Sprintf("%d", params.Port),
	}
	configYAML, err := yaml.Marshal(config)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	// 保存

	var configFile = Tea.ConfigFile("api_db.yaml")
	err = os.WriteFile(configFile, configYAML, 0666)
	if err != nil {
		this.Fail("保存配置失败：" + err.Error())
		return
	}

	// TODO 思考是否让本地的API节点生效

	this.Success()
}

func (this *UpdateAction) parseOldConfig(data []byte) {
	var config = &dbs.Config{}
	err := yaml.Unmarshal(data, config)
	if err != nil {
		this.Show()
		return
	}

	if config.DBs == nil {
		this.Show()
		return
	}

	var dbConfig *dbs.DBConfig
	for _, db := range config.DBs {
		dbConfig = db
		break
	}
	if dbConfig == nil {
		this.Data["dbConfig"] = maps.Map{
			"host":     "",
			"port":     "",
			"username": "",
			"password": "",
			"database": "",
		}
		this.Show()
		return
	}

	var dsn = dbConfig.Dsn
	cfg, err := mysql.ParseDSN(dsn)
	if err != nil {
		this.Data["dbConfig"] = maps.Map{
			"host":     "",
			"port":     "",
			"username": "",
			"password": "",
			"database": "",
		}
		this.Show()
		return
	}

	var host = cfg.Addr
	var port = "3306"
	var index = strings.LastIndex(cfg.Addr, ":")
	if index > 0 {
		host = cfg.Addr[:index]
		port = cfg.Addr[index+1:]
	}

	this.Data["dbConfig"] = maps.Map{
		"host":     host,
		"port":     port,
		"username": cfg.User,
		"password": cfg.Passwd,
		"database": cfg.DBName,
	}

	this.Show()
}
