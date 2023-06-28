//go:build !plus

package settingutils

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/configloaders"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/helpers"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
	"github.com/iwind/TeaGo/actions"
)

type Helper struct {
	helpers.LangHelper

	tab string
}

func NewHelper(tab string) *Helper {
	return &Helper{
		tab: tab,
	}
}

func (this *Helper) BeforeAction(actionPtr actions.ActionWrapper) (goNext bool) {
	goNext = true

	var action = actionPtr.Object()

	// 左侧菜单
	action.Data["teaMenu"] = "settings"
	action.Data["teaSubMenu"] = "basic"

	// 标签栏
	var tabbar = actionutils.NewTabbar()
	var session = action.Session()
	var adminId = session.GetInt64("adminId")
	if configloaders.AllowModule(adminId, configloaders.AdminModuleCodeSetting) {
		tabbar.Add(this.Lang(actionPtr, codes.AdminSettingsTabAdminServer), "", "/settings/server", "", this.tab == "server")
		tabbar.Add(this.Lang(actionPtr, codes.AdminSettingsTabAdminUI), "", "/settings/ui", "", this.tab == "ui")
		tabbar.Add(this.Lang(actionPtr, codes.AdminSettingsTabAdminSecuritySettings), "", "/settings/security", "", this.tab == "security")
		tabbar.Add(this.Lang(actionPtr, codes.AdminSettingsTabUpdates), "", "/settings/updates", "", this.tab == "updates")
	}
	tabbar.Add(this.Lang(actionPtr, codes.AdminSettingsTabProfile), "", "/settings/profile", "", this.tab == "profile")
	tabbar.Add(this.Lang(actionPtr, codes.AdminSettingsTabLogin), "", "/settings/login", "", this.tab == "login")
	actionutils.SetTabbar(actionPtr, tabbar)

	return
}
