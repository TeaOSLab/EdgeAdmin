package waf

import "github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"

type ExportAction struct {
	actionutils.ParentAction
}

func (this *ExportAction) Init() {
	this.Nav("", "", "export")
}

func (this *ExportAction) RunGet(params struct{}) {
	this.Show()
}
