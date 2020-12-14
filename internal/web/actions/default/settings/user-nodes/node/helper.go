package node

import (
	"github.com/iwind/TeaGo/actions"
	"net/http"
)

type Helper struct {
}

func NewHelper() *Helper {
	return &Helper{}
}

func (this *Helper) BeforeAction(action *actions.ActionObject) (goNext bool) {
	if action.Request.Method != http.MethodGet {
		return true
	}

	return true
}
