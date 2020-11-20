package profile

import (
	"fmt"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/dbs"
	"github.com/iwind/TeaGo/maps"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"net/url"
	"path/filepath"
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
	dsn = regexp.MustCompile(`tcp\((.+)\)`).ReplaceAllString(dsn, "$1")
	dsnURL, err := url.Parse("mysql://" + dsn)
	if err != nil {
		this.Show()
		return
	}

	host := dsnURL.Host
	port := "3306"
	index := strings.LastIndex(dsnURL.Host, ":")
	if index > 0 {
		host = dsnURL.Host[:index]
		port = dsnURL.Host[index+1:]
	}

	password, _ := dsnURL.User.Password()
	this.Data["dbConfig"] = maps.Map{
		"host":     host,
		"port":     port,
		"username": dsnURL.User.Username(),
		"password": password,
		"database": filepath.Base(dsnURL.Path),
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
		Match(`^[\w\.-]+$`, "主机地址中不能包含特殊字符").
		Field("port", params.Port).
		Gt(0, "端口需要大于0").
		Lt(65535, "端口需要小于65535").
		Field("database", params.Database).
		Require("请输入数据库名称").
		Match(`^[\w\.-]+$`, "数据库名称中不能包含特殊字符").
		Field("username", params.Username).
		Require("请输入连接数据库的用户名").
		Match(`^[\w\.-]+$`, "用户名中不能包含特殊字符")

	if len(params.Password) > 0 {
		params.Must.
			Field("password", params.Password).
			Match(`^[\w\.-]+$`, "密码中不能包含特殊字符")
	}

	// 保存
	dsn := params.Username + ":" + params.Password + "@tcp(" + params.Host + ":" + fmt.Sprintf("%d", params.Port) + ")/" + params.Database

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
