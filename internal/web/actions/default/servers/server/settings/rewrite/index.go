package rewrite

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.Nav("", "setting", "index")
	this.SecondMenu("rewrite")
}

func (this *IndexAction) RunGet(params struct {
	ServerId int64
}) {
	// TODO

	this.Show()
}
