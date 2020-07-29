package nodes

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/iwind/TeaGo/actions"
)

type Helper struct {
}

func (this *Helper) BeforeAction(action *actions.ActionObject) {
	action.Data["teaMenu"] = "nodes"

	selectedTabbar, _ := action.Data["mainTab"]

	tabbar := actionutils.NewTabbar()
	tabbar.Add("节点管理", "", "/nodes", "", selectedTabbar == "node")
	tabbar.Add("认证管理", "", "/nodes/grants", "", selectedTabbar == "grant")
	actionutils.SetTabbar(action, tabbar)
}
