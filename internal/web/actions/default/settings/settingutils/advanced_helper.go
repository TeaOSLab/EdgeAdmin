package settingutils

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/configloaders"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/iwind/TeaGo/actions"
)

type AdvancedHelper struct {
	tab string
}

func NewAdvancedHelper(tab string) *AdvancedHelper {
	return &AdvancedHelper{
		tab: tab,
	}
}

func (this *AdvancedHelper) BeforeAction(actionPtr actions.ActionWrapper) (goNext bool) {
	goNext = true

	action := actionPtr.Object()

	// 左侧菜单
	action.Data["teaMenu"] = "settings"
	action.Data["teaSubMenu"] = "advanced"

	// 标签栏
	tabbar := actionutils.NewTabbar()
	var session = action.Session()
	var adminId = session.GetInt64("adminId")
	if configloaders.AllowModule(adminId, configloaders.AdminModuleCodeSetting) {
		tabbar.Add("数据库", "", "/settings/database", "", this.tab == "database")
		tabbar.Add("API节点", "", "/api", "", this.tab == "apiNodes")
		tabbar.Add("用户节点", "", "/settings/userNodes", "", this.tab == "userNodes")
		tabbar.Add("日志数据库", "", "/db", "", this.tab == "dbNodes")
	}
	actionutils.SetTabbar(actionPtr, tabbar)

	return
}
