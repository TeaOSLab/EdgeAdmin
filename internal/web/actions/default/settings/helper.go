package settings

import "github.com/iwind/TeaGo/actions"

type Helper struct {
}

func NewHelper() *Helper {
	return &Helper{}
}

func (this *Helper) BeforeAction(action *actions.ActionObject) {
	action.Data["teaMenu"] = "settings"
}
