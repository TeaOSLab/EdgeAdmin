package issues

import "github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"

type UpdateServerPopupAction struct {
	actionutils.ParentAction
}

func (this *UpdateServerPopupAction) Init() {
	this.Nav("", "", "")
}

func (this *UpdateServerPopupAction) RunGet(params struct{}) {
	this.Show()
}
