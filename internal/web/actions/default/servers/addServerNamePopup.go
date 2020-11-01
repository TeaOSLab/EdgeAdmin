package servers

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
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
		params.Must.
			Field("serverName", params.ServerName).
			Require("请输入域名")
		this.Data["serverName"] = maps.Map{
			"name": params.ServerName,
			"type": "full",
		}
	} else if params.Mode == "multiple" {
		if len(params.ServerNames) == 0 {
			this.FailField("serverNames", "请输入至少域名")
		}

		serverNames := []string{}
		for _, line := range strings.Split(params.ServerNames, "\n") {
			line := strings.TrimSpace(line)
			line = regexp.MustCompile(`\s+`).ReplaceAllString(line, "")
			if len(line) == 0 {
				continue
			}
			serverNames = append(serverNames, line)
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
