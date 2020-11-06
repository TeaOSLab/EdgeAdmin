package ipadmin

import "github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"

type ListsAction struct {
	actionutils.ParentAction
}

func (this *ListsAction) Init() {
	this.Nav("", "", "ipadmin")
}

func (this *ListsAction) RunGet(params struct{}) {
	this.Data["subMenuItem"] = "list"

	this.Show()
}
