package waf

import "github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"

type LogAction struct {
	actionutils.ParentAction
}

func (this *LogAction) Init() {
	this.Nav("", "", "")
}

func (this *LogAction) RunGet(params struct{}) {
	this.Show()
}
