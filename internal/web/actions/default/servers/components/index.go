package components

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.Nav("", "component", "index")
}

func (this *IndexAction) RunGet(params struct{}) {

	this.Show()
}
