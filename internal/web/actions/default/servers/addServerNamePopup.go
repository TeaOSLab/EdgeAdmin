package servers

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
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
	ServerName string

	Must *actions.Must
}) {
	params.Must.
		Field("serverName", params.ServerName).
		Require("请输入域名")

	this.Data["serverName"] = maps.Map{
		"name": params.ServerName,
		"type": "full",
	}
	this.Success()
}
