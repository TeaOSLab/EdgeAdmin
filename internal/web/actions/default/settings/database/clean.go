package profile

import "github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"

type CleanAction struct {
	actionutils.ParentAction
}

func (this *CleanAction) Init() {
	this.Nav("", "", "clean")
}

func (this *CleanAction) RunGet(params struct{}) {
	this.Show()
}
