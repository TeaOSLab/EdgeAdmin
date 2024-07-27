package node

import (
	"net/http"

	"github.com/iwind/TeaGo/actions"
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
