package messages

import "github.com/iwind/TeaGo/actions"

type Helper struct {
}

func (this *Helper) BeforeAction(actionPtr actions.ActionWrapper) {
	action := actionPtr.Object()
	action.Data["teaMenu"] = "message"
}