package database

import (
	"fmt"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/configutils"
	"github.com/go-sql-driver/mysql"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/dbs"
	"github.com/iwind/TeaGo/maps"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"net"
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
	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		return
	}

	config := &dbs.Config{}
	err = yaml.Unmarshal(data, config)
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

	dsn := dbConfig.Dsn
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

	host := cfg.Addr
	port := "3306"
	index := strings.LastIndex(cfg.Addr, ":")
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

func (this *UpdateAction) RunPost(params struct {
	Host     string
	Port     int32
	Database string
	Username string
	Password string

	Must *actions.Must
}) {
	defer this.CreateLogInfo("修改API节点数据库设置")

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

	// 保存
	dsn := params.Username + ":" + params.Password + "@tcp(" + configutils.QuoteIP(params.Host) + ":" + fmt.Sprintf("%d", params.Port) + ")/" + params.Database

	configFile := Tea.ConfigFile("api_db.yaml")
	template := `default:
  db: "prod"
  prefix: ""

dbs:
  prod:
    driver: "mysql"
    dsn: "` + dsn + `?charset=utf8mb4&timeout=30s"
    prefix: "edge"
    models:
      package: internal/web/models
`
	err := ioutil.WriteFile(configFile, []byte(template), 0666)
	if err != nil {
		this.Fail("保存配置失败：" + err.Error())
	}

	// TODO 让本地的节点生效

	this.Success()
}
