package cache

import "github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"

type PurgeAction struct {
	actionutils.ParentAction
}

func (this *PurgeAction) Init() {
	this.Nav("", "", "purge")
}

func (this *PurgeAction) RunGet(params struct{}) {
	this.Show()
}
