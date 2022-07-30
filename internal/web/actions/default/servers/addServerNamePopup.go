package servers

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
	"net/url"
	"regexp"
	"strings"
)

type AddServerNamePopupAction struct {
	actionutils.ParentAction
}

func (this *AddServerNamePopupAction) Init() {
	this.Nav("", "", "")
}

func (this *AddServerNamePopupAction) RunGet(params struct{}) {
	this.Show()
}

func (this *AddServerNamePopupAction) RunPost(params struct {
	Mode string

	ServerName  string
	ServerNames string

	Must *actions.Must
}) {
	if params.Mode == "single" {
		var serverName = strings.ToLower(params.ServerName)

		// 去除空格
		serverName = regexp.MustCompile(`\s+`).ReplaceAllString(serverName, "")

		// 处理URL
		if regexp.MustCompile(`^(?i)(http|https|ftp)://`).MatchString(serverName) {
			u, err := url.Parse(serverName)
			if err == nil && len(u.Host) > 0 {
				serverName = u.Host
			}
		}

		params.Must.
			Field("serverName", serverName).
			Require("请输入域名")

		this.Data["serverName"] = maps.Map{
			"name": serverName,
			"type": "full",
		}
	} else if params.Mode == "multiple" {
		if len(params.ServerNames) == 0 {
			this.FailField("serverNames", "请输入至少域名")
		}

		var serverNames = []string{}
		for _, line := range strings.Split(params.ServerNames, "\n") {
			var serverName = strings.TrimSpace(line)
			serverName = regexp.MustCompile(`\s+`).ReplaceAllString(serverName, "")
			if len(serverName) == 0 {
				continue
			}

			// 处理URL
			if regexp.MustCompile(`^(?i)(http|https|ftp)://`).MatchString(serverName) {
				u, err := url.Parse(serverName)
				if err == nil && len(u.Host) > 0 {
					serverName = u.Host
				}
			}

			// 转成小写
			serverName = strings.ToLower(serverName)

			serverNames = append(serverNames, serverName)
		}
		this.Data["serverName"] = maps.Map{
			"name":     "",
			"type":     "full",
			"subNames": serverNames,
		}
	} else {
		this.Fail("错误的mode参数")
	}

	this.Success()
}
