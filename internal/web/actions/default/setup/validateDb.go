package setup

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/configutils"
	_ "github.com/go-sql-driver/mysql"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/dbs"
	"github.com/iwind/TeaGo/maps"
	stringutil "github.com/iwind/TeaGo/utils/string"
	"net"
	"net/url"
	"regexp"
	"strings"
)

type ValidateDbAction struct {
	actionutils.ParentAction
}

func (this *ValidateDbAction) RunPost(params struct {
	Host              string
	Port              string
	Database          string
	Username          string
	Password          string
	AccessLogKeepDays int

	Must *actions.Must
}) {
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
		Require("请输入端口").
		Match(`^\d+$`, "端口中只能包含数字").
		Field("database", params.Database).
		Require("请输入数据库名称").
		Match(`^[\w\.-]+$`, "数据库名称中不能包含特殊字符").
		Field("username", params.Username).
		Require("请输入连接数据库的用户名").
		Match(`^[\w\.-]+$`, "用户名中不能包含特殊字符")

	// 测试连接
	db, err := dbs.NewInstanceFromConfig(&dbs.DBConfig{
		Driver: "mysql",
		Dsn:    url.QueryEscape(params.Username) + ":" + url.QueryEscape(params.Password) + "@tcp(" + configutils.QuoteIP(params.Host) + ":" + params.Port + ")/" + params.Database,
		Prefix: "",
	})
	if err != nil {
		this.Fail("数据库信息错误：" + err.Error())
	}

	defer func() {
		_ = db.Close()
	}()

	err = db.Raw().Ping()
	if err != nil {
		// 是否是数据库不存在
		if strings.Contains(err.Error(), "Error 1049") {
			db, err := dbs.NewInstanceFromConfig(&dbs.DBConfig{
				Driver: "mysql",
				Dsn:    params.Username + ":" + params.Password + "@tcp(" + configutils.QuoteIP(params.Host) + ":" + params.Port + ")/",
				Prefix: "",
			})
			if err != nil {
				this.Fail("尝试创建数据库失败：" + err.Error())
				return
			}

			_, err = db.Exec("CREATE DATABASE `" + params.Database + "`")
			if err != nil {
				this.Fail("尝试创建数据库失败：" + err.Error())
				return
			}
		} else {
			if strings.Contains(err.Error(), "Error 1044:") {
				this.Fail("无法连接到数据库，权限检查失败：" + err.Error())
				return
			}
			this.Fail("无法连接到数据库，请检查配置：" + err.Error())
			return
		}
	}

	// 检查权限
	// edgeTest表名需要根据表结构的变更而变更，防止升级时冲突
	var testTable = "edgeTest1"
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS `" + testTable + "` (\n  `id` int(11) NOT NULL AUTO_INCREMENT,\n  PRIMARY KEY (`id`)\n) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;")
	if err != nil {
		this.Fail("当前连接的用户无法创建新表，请检查CREATE权限设置：" + err.Error())
	}

	_, err = db.Exec("ALTER TABLE `" + testTable + "` CHANGE `id` `id` int(11) NOT NULL AUTO_INCREMENT")
	if err != nil {
		this.Fail("当前连接的用户无法修改表结构，请检查ALTER权限设置：" + err.Error())
	}

	// 删除edgeTest，忽略可能的错误，因为我们不需要DROP权限
	_, _ = db.Exec("DROP TABLE `" + testTable + "`")

	// 检查数据库版本
	one, err := db.FindOne("SELECT VERSION() AS v")
	if err != nil {
		this.Fail("检查数据库版本时出错：" + err.Error())
	}
	if one == nil {
		this.Fail("检查数据库版本时出错：无法获取数据库版本")
	}
	version := one.GetString("v")
	if stringutil.VersionCompare(version, "5.7.8") < 0 {
		this.Fail("数据库版本至少在v5.7.8以上，你现在使用的是v" + version)
	}

	this.Data["db"] = maps.Map{
		"host":              params.Host,
		"port":              params.Port,
		"database":          params.Database,
		"username":          params.Username,
		"password":          params.Password,
		"passwordMask":      strings.Repeat("*", len(params.Password)),
		"accessLogKeepDays": params.AccessLogKeepDays,
	}

	this.Success()
}
