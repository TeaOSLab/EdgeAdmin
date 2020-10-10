package server

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"strconv"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.Nav("", "index", "index")
	this.SecondMenu("index")
}

func (this *IndexAction) RunGet(params struct {
	ServerId int64
}) {
	// TODO 等看板实现后，需要跳转到看板
	this.RedirectURL("/servers/server/log?serverId=" + strconv.FormatInt(params.ServerId, 10))
}
