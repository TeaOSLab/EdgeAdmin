package waf

import "github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"

type GroupsAction struct {
	actionutils.ParentAction
}

func (this *GroupsAction) Init() {
	this.Nav("", "", "")
}

func (this *GroupsAction) RunGet(params struct{}) {
	this.Show()
}
