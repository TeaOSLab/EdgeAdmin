package clusters

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
	action.Data["teaMenu"] = "clusters"

	selectedTabbar, _ := action.Data["mainTab"]

	tabbar := actionutils.NewTabbar()
	tabbar.Add("集群", "", "/clusters", "", selectedTabbar == "cluster")
	actionutils.SetTabbar(action, tabbar)
}
