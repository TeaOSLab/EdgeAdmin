package server

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.Nav("", "setting", "index")
	this.SecondMenu("http")
}

func (this *IndexAction) RunGet(params struct {
	ServerId int64
}) {
	this.Show()
}
