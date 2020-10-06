package waf

import "github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"

type ImportAction struct {
	actionutils.ParentAction
}

func (this *ImportAction) Init() {
	this.Nav("", "", "")
}

func (this *ImportAction) RunGet(params struct{}) {
	this.Show()
}
