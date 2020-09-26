package websocket

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/iwind/TeaGo/actions"
)

// 添加来源域
type CreateOriginAction struct {
	actionutils.ParentAction
}

func (this *CreateOriginAction) Init() {
	this.Nav("", "", "")
}

func (this *CreateOriginAction) RunGet(params struct{}) {
	this.Show()
}

func (this *CreateOriginAction) RunPost(params struct {
	Origin string

	Must *actions.Must
}) {
	params.Must.
		Field("origin", params.Origin).
		Require("请输入域名")

	this.Data["origin"] = params.Origin
	this.Success()
}
