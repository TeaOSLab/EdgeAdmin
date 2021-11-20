package settingutils

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/configloaders"
	teaconst "github.com/TeaOSLab/EdgeAdmin/internal/const"
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
		if teaconst.IsPlus {
			tabbar.Add("用户节点", "", "/settings/userNodes", "", this.tab == "userNodes")
		}
		tabbar.Add("日志数据库", "", "/db", "", this.tab == "dbNodes")
		if teaconst.IsPlus {
			tabbar.Add("监控节点", "", "/settings/monitorNodes", "", this.tab == "monitorNodes")
		}
		tabbar.Add("迁移", "", "/settings/transfer", "", this.tab == "transfer")
		if teaconst.BuildPlus {
			tabbar.Add("商业版认证", "", "/settings/authority", "", this.tab == "authority")
		}

		//tabbar.Add("备份", "", "/settings/backup", "", this.tab == "backup")
	}
	actionutils.SetTabbar(actionPtr, tabbar)

	return
}
