package groups

import (
	"github.com/iwind/TeaGo/actions"
	"net/http"
)

type Helper struct {
}

func NewHelper() *Helper {
	return &Helper{}
}

func (this *Helper) BeforeAction(action *actions.ActionObject) {
	if action.Request.Method != http.MethodGet {
		return
	}

	action.Data["teaMenu"] = "servers"
	action.Data["mainTab"] = "component"
	action.Data["secondMenuItem"] = "group"
}
