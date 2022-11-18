// Copyright 2021 Liuxiangchao iwind.liu@gmail.com. All rights reserved.

package setup

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/configutils"
	"github.com/iwind/TeaGo/dbs"
	"github.com/iwind/TeaGo/maps"
	"net"
	"os"
	"strings"
	"time"
)

// DetectDBAction 尝试从本地服务器中发现MySQL
type DetectDBAction struct {
	actionutils.ParentAction
}

func (this *DetectDBAction) RunPost(params struct{}) {
	var localHost = ""
	var localPort = ""
	var localUsername = ""
	var localPassword = ""

	// 本地的3306端口是否可以连接
	conn, err := net.DialTimeout("tcp", "127.0.0.1:3306", 3*time.Second)
	if err == nil {
		_ = conn.Close()
		localHost = "127.0.0.1"
		localPort = "3306"

		var username = "root"
		var passwords = []string{"", "123456", "654321", "Aa_123456", "111111"}

		// 使用 foolish-mysql 安装的MySQL
		localGeneratedPasswordData, err := os.ReadFile("/usr/local/mysql/generated-password.txt")
		if err == nil {
			var localGeneratedPassword = strings.TrimSpace(string(localGeneratedPasswordData))
			if len(localGeneratedPassword) > 0 {
				passwords = append(passwords, localGeneratedPassword)
			}
		}

		for _, pass := range passwords {
			db, err := dbs.NewInstanceFromConfig(&dbs.DBConfig{
				Driver: "mysql",
				Dsn:    username + ":" + pass + "@tcp(" + configutils.QuoteIP(localHost) + ":" + localPort + ")/edges",
				Prefix: "",
			})
			if err == nil {
				err = db.Raw().Ping()
				_ = db.Close()

				if err == nil || strings.Contains(err.Error(), "Error 1049") {
					localUsername = username
					localPassword = pass
					break
				}
			}
		}
	}

	this.Data["localDB"] = maps.Map{
		"host":     localHost,
		"port":     localPort,
		"username": localUsername,
		"password": localPassword,
	}

	this.Success()
}
