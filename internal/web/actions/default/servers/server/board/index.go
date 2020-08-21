package board

import "github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.Nav("", "board", "")
}

func (this *IndexAction) RunGet(params struct{}) {
	this.Show()
}
