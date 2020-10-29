package groups

import "github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"

type OptionsAction struct {
	actionutils.ParentAction
}

func (this *OptionsAction) RunPost(params struct{}) {
	this.Success()
}
