package ipadmin

import "github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"

type ProvincesAction struct {
	actionutils.ParentAction
}

func (this *ProvincesAction) Init() {
	this.Nav("", "", "ipadmin")
}

func (this *ProvincesAction) RunGet(params struct{}) {
	this.Data["subMenuItem"] = "province"
	this.Show()
}
