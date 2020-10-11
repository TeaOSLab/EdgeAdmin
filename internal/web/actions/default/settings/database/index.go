package profile

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/go-yaml/yaml"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/dbs"
	"github.com/iwind/TeaGo/maps"
	"io/ioutil"
	"net/url"
	"path/filepath"
	"regexp"
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

	configFile := Tea.ConfigFile("api_db.yaml")
	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		this.Data["error"] = "read config file failed: api_db.yaml: " + err.Error()
		this.Show()
		return
	}

	config := &dbs.Config{}
	err = yaml.Unmarshal(data, config)
	if err != nil {
		this.Data["error"] = "parse config file failed: api_db.yaml: " + err.Error()
		this.Show()
		return
	}

	if config.DBs == nil {
		this.Data["error"] = "can not find valid database config: api_db.yaml"
		this.Show()
		return
	}

	dbConfig, ok := config.DBs[config.Default.DB]
	if !ok {
		this.Data["error"] = "can not find valid database config: api_db.yaml"
		this.Show()
		return
	}
	dsn := dbConfig.Dsn
	dsn = regexp.MustCompile(`tcp\((.+)\)`).ReplaceAllString(dsn, "$1")
	dsnURL, err := url.Parse("mysql://" + dsn)
	if err != nil {
		this.Data["error"] = "parse dsn failed: " + err.Error()
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
	if len(password) > 0 {
		password = strings.Repeat("*", len(password))
	}
	this.Data["dbConfig"] = maps.Map{
		"host":     host,
		"port":     port,
		"username": dsnURL.User.Username(),
		"password": password,
		"database": filepath.Base(dsnURL.Path),
	}

	// TODO 测试连接

	this.Show()
}
