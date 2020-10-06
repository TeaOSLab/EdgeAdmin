package waf

import "github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"

type SetsAction struct {
	actionutils.ParentAction
}

func (this *SetsAction) Init() {
	this.Nav("", "", "")
}

func (this *SetsAction) RunGet(params struct{}) {
	this.Show()
}
