package cache

import "github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"

type PreheatAction struct {
	actionutils.ParentAction
}

func (this *PreheatAction) Init() {
	this.Nav("", "", "preheat")
}

func (this *PreheatAction) RunGet(params struct{}) {
	this.Show()
}
