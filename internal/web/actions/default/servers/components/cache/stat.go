package cache

import "github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"

type StatAction struct {
	actionutils.ParentAction
}

func (this *StatAction) Init() {
	this.Nav("", "", "stat")
}

func (this *StatAction) RunGet(params struct{}) {
	this.Show()
}
