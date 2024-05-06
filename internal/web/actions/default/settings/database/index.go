package database

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/configs"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/go-sql-driver/mysql"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/dbs"
	"github.com/iwind/TeaGo/maps"
	"gopkg.in/yaml.v3"
	"net"
	"os"
	"strings"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.Nav("", "", "index")
}

func (this *IndexAction) RunGet(params struct{}) {
	this.Data["error"] = ""

	var configFile = Tea.ConfigFile("api_db.yaml")
	data, err := os.ReadFile(configFile)
	if err != nil {
		this.Data["error"] = "read config file failed: api_db.yaml: " + err.Error()
		this.Show()
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

func (this *IndexAction) parseOldConfig(data []byte) {
	var config = &dbs.Config{}
	err := yaml.Unmarshal(data, config)
	if err != nil {
		this.Data["error"] = "parse config file failed: api_db.yaml: " + err.Error()
		this.Show()
		return
	}

	if config.DBs == nil {
		this.Data["error"] = "no database configured in config file: api_db.yaml"
		this.Show()
		return
	}

	var dbConfig *dbs.DBConfig
	for _, db := range config.DBs {
		dbConfig = db
		break
	}
	if dbConfig == nil {
		this.Data["error"] = "no database configured in config file: api_db.yaml"
		this.Show()
		return
	}
	var dsn = dbConfig.Dsn
	cfg, err := mysql.ParseDSN(dsn)
	if err != nil {
		this.Data["error"] = "parse dsn error: " + err.Error()
		this.Show()
		return
	}

	var host = cfg.Addr
	var port = "3306"
	var index = strings.LastIndex(host, ":")
	if index > 0 {
		port = host[index+1:]
		host = host[:index]
	}

	var password = cfg.Passwd
	if len(password) > 0 {
		password = strings.Repeat("*", len(password))
	}

	this.Data["dbConfig"] = maps.Map{
		"host":     host,
		"port":     port,
		"username": cfg.User,
		"password": password,
		"database": cfg.DBName,
	}

	// TODO 测试连接

	this.Show()
}
