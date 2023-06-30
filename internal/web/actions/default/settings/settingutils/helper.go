//go:build !plus

package settingutils

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/configloaders"
	teaconst "github.com/TeaOSLab/EdgeAdmin/internal/const"
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
	var adminId = session.GetInt64(teaconst.SessionAdminId)
	if configloaders.AllowModule(adminId, configloaders.AdminModuleCodeSetting) {
		tabbar.Add(this.Lang(actionPtr, codes.AdminSetting_TabAdminServer), "", "/settings/server", "", this.tab == "server")
		tabbar.Add(this.Lang(actionPtr, codes.AdminSetting_TabAdminUI), "", "/settings/ui", "", this.tab == "ui")
		tabbar.Add(this.Lang(actionPtr, codes.AdminSetting_TabAdminSecuritySettings), "", "/settings/security", "", this.tab == "security")
		tabbar.Add(this.Lang(actionPtr, codes.AdminSetting_TabUpdates), "", "/settings/updates", "", this.tab == "updates")
	}
	tabbar.Add(this.Lang(actionPtr, codes.AdminSetting_TabProfile), "", "/settings/profile", "", this.tab == "profile")
	tabbar.Add(this.Lang(actionPtr, codes.AdminSetting_TabLogin), "", "/settings/login", "", this.tab == "login")
	actionutils.SetTabbar(actionPtr, tabbar)

	return
}
