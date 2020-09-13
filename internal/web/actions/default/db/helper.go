package db

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/iwind/TeaGo/actions"
	"net/http"
)

type Helper struct {
}

func (this *Helper) BeforeAction(action *actions.ActionObject) {
	if action.Request.Method != http.MethodGet {
		return
	}

	action.Data["teaMenu"] = "db"

	selectedTabbar, _ := action.Data["mainTab"]

	tabbar := actionutils.NewTabbar()
	tabbar.Add("数据库节点", "", "/db", "", selectedTabbar == "db")
	actionutils.SetTabbar(action, tabbar)
}
