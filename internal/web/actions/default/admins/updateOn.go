package admins

import "github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"

type UpdateOnAction struct {
	actionutils.ParentAction
}

func (this *UpdateOnAction) RunPost(params struct{}) {
	this.Success()
}
