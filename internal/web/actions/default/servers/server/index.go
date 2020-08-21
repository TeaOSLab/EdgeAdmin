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
	this.RedirectURL("/servers/server/board?serverId=" + strconv.FormatInt(params.ServerId, 10))
}
