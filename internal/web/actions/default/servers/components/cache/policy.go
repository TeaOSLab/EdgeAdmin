package cache

import "github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"

type PolicyAction struct {
	actionutils.ParentAction
}

func (this *PolicyAction) Init() {
	this.Nav("", "", "")
}

func (this *PolicyAction) RunGet(params struct{}) {
	this.Show()
}
