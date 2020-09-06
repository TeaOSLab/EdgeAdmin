package api

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/iwind/TeaGo/actions"
)

type Helper struct {
}

func NewHelper() *Helper {
	return &Helper{}
}

func (this *Helper) BeforeAction(action *actions.ActionObject) {
	action.Data["teaMenu"] = "api"

	selectedTabbar, _ := action.Data["mainTab"]

	tabbar := actionutils.NewTabbar()
	tabbar.Add("API节点", "", "/api", "", selectedTabbar == "node")
	actionutils.SetTabbar(action, tabbar)
}
