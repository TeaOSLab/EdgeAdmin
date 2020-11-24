package acme

import "github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"

type CreateAction struct {
	actionutils.ParentAction
}

func (this *CreateAction) Init() {
	this.Nav("", "", "")
}

func (this *CreateAction) RunGet(params struct{}) {
	this.Show()
}
