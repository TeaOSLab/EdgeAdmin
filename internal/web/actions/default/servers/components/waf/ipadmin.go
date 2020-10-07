package waf

import "github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"

type IpadminAction struct {
	actionutils.ParentAction
}

func (this *IpadminAction) Init() {
	this.Nav("", "", "ipadmin")
}

func (this *IpadminAction) RunGet(params struct{}) {
	this.Show()
}
