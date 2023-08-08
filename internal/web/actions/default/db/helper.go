package db

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/helpers"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
	"github.com/iwind/TeaGo/actions"
	"net/http"
)

type Helper struct {
	helpers.LangHelper
}

func (this *Helper) BeforeAction(action *actions.ActionObject) {
	if action.Request.Method != http.MethodGet {
		return
	}

	action.Data["teaMenu"] = "db"

	var selectedTabbar = action.Data["mainTab"]

	var tabbar = actionutils.NewTabbar()
	tabbar.Add(this.Lang(action, codes.DBNode_TabNodes), "", "/db", "", selectedTabbar == "db")
	actionutils.SetTabbar(action, tabbar)
}
